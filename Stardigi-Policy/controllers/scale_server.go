package controllers

import (
	"Stardigi-Policy/config"
	"Stardigi-Policy/model/db"
	"Stardigi-Policy/model/httpc"
	"Stardigi-Policy/routers"
	"Stardigi-Policy/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	// "sync"
	"time"
)

var (
	LastScale   *utils.SyncMap
	LastCollect *utils.SyncMap

	// 用于判断当前正扩缩容的规则是否完成，或者是否移除
	scaling *utils.SyncMap
)

func init() {
	LastScale = utils.NewSyncMap()
	LastCollect = utils.NewSyncMap()
	scaling = utils.NewSyncMap()
}

// // 定义全局变量存储每条记录上次扩容时间
// var lastScale = make(map[string]int64, 0)
// var lastCollect = make(map[string]int64, 0)

// 设置变量countMatch== collect_period时执行扩缩容
// 1. 获取规则 2. 根据规则去获取指标信息 3.
func SRun(policy *config.PolicyConfig) {

	signals := make(chan os.Signal, 0)
	signal.Notify(signals, os.Interrupt, os.Kill)
	// 启动协程老化scaling中的数据
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				scaling.Age()
			case <-signals:
				fmt.Println("由于系统执行了中断命令，扩缩容程序退出...")
				return
			}
		}
	}()

	for {
		select {
		case ruleUp := <-RuleUpChan:
			// 执行扩容操作
			//根据app_id和marathon_name从Prometheus获取监控数据
			fmt.Println("==查询扩容规则完成==", ruleUp)
			scaleJobs(ruleUp, policy)
		case ruleDown := <-RuleDownChan:
			// 执行缩容操作
			// 问题：速度太快,导致一条规则多次加入扩缩容的策略
			// 考虑添加控制语句要求上次规则必须使用完毕
			// fmt.Println(ruleDown)
			fmt.Println("==查询缩容规则完成==", ruleDown)
			scaleJobs(ruleDown, policy)

		case <-signals:
			fmt.Println("由于系统执行了中断命令，扩缩容程序退出...")
			return
		}
	}
}

func scaleJobs(rules []db.AppScaleRule, policy *config.PolicyConfig) {

	for _, rule := range rules {
		ruleFlag := utils.StructJoin(&rule)
		//查看map是否存在这条规则
		if v, ok := scaling.Get(ruleFlag); ok {
			// 说明规则正在扩缩容
			if v.(bool) {
				continue
			}
		}

		scaling.Set(ruleFlag, true)

		//每条规则启动线程扩缩容
		// 每一条规则
		//判断有没有达到扩缩容冷切时间
		go func(rule db.AppScaleRule) {

			appInfo := utils.StringJoin(rule.MarathonName, rule.AppId, rule.ScaleType)
			var lastTimeScale int64 = 0
			currentTimeScale := time.Now().Unix()
			if v, ok := LastScale.Get(appInfo); ok {
				lastTimeScale = v.(int64)
			}
			// 数据库中以分为单位
			if currentTimeScale-lastTimeScale < int64(rule.ColdTime*60) {
				scaling.Set(ruleFlag, false)
				return
			}

			var exit bool = false
			// 判断有没有达到收集冷切时间

			memChan := make(chan bool, 3)
			cpuChan := make(chan bool, 3)
			haChan := make(chan bool, 3)
			threadChan := make(chan bool, 3)
			counter := make(map[string]chan bool, 0)
			counter["memory"] = memChan
			counter["cpu"] = cpuChan
			counter["ha_queue"] = haChan
			counter["thread"] = threadChan

			// // 启动检查程序
			// // 问题：启动了多个次协程

			go func(rule db.AppScaleRule, counter map[string]chan bool) {
				//标识扩缩容完成
				finScaled := make(chan bool, 1)
				matched := rule.ContinuePeriod
				for {

					select {
					case <-time.After(5 * time.Second):
						if len(counter["memory"]) == matched || len(counter["cpu"]) == matched || len(counter["ha_queue"]) == matched || len(counter["thread"]) == matched {
							// 执行扩缩容
							// 清空数据，防止再次进入

							SetChanNull(counter)

							scaleJob(policy, rule, finScaled)
						}
					case <-finScaled:
						// 记录扩缩容时间
						appscaledInfo := utils.StringJoin(rule.MarathonName, rule.AppId, rule.ScaleType)
						LastScale.Set(appscaledInfo, time.Now().Unix())
						exit = true
						ruleFlag := utils.StructJoin(&rule)
						scaling.Set(ruleFlag, false)
						return
					}
				}
			}(rule, counter)

			for {

				if exit {
					return
				}

				//先判断有收集冷切时间有没有达到
				//暂时方案，考虑该冷切时间的判断放到metric_server.go中
				var lastTimeCollect int64 = 0
				currentTimeCollect := time.Now().Unix()
				if v, ok := LastCollect.Get(appInfo); ok {
					lastTimeCollect = v.(int64)
				}

				if currentTimeCollect-lastTimeCollect < int64(rule.CollectPeriod*60) {
					continue
				}
				//已达到时间
				if strings.EqualFold(rule.ScaleType, "up") {

					RulesUp <- rule
					// Metrics为无缓存通道
					metrics := <-MetricsUp
					fmt.Println("==从prometheus查询到指标信息==", metrics)
					LastCollect.Set(appInfo, time.Now().Unix())
					//
					matchJobs(rule, metrics[appInfo], counter)
				} else if strings.EqualFold(rule.ScaleType, "down") {

					RulesDown <- rule
					// Metrics为无缓存通道
					metrics := <-MetricsDown
					LastCollect.Set(appInfo, time.Now().Unix())
					matchJobs(rule, metrics[appInfo], counter)
				}
			}
		}(rule)
	}
}

/**
 * 执行扩容操作
 */
func scaleJob(policy *config.PolicyConfig, rule db.AppScaleRule, finScaled chan bool) {
	/***************获取marathon地址***************/
	marathonEndpoint := make(chan string, 1)
	marathon_name := os.Getenv("marathon_name")
	marathon_url := os.Getenv("marathon_url")
	marathon_endpoint := utils.StringJoin(marathon_url, "/", marathon_name)
	marathonEndpoint <- marathon_endpoint
	// "10.254.9.55:5050;10.254.9.56:5050;10.254.9.57:5050"
	// mesosUrls := strings.Split(policy.Mesoss.MesosUrls, ";")

	// // 启动异步请求，获取marathon地址
	// finMesos := make(chan bool, 1)
	// var once sync.Once
	// for _, v := range mesosUrls {

	// 	go func(v string, marathonUrl chan string) {
	// 		host := utils.StringJoin("http://", v)
	// 		starRequestGen := httpc.NewStarRequestGen(host, routers.MesosRouter)
	// 		responseMesosChan, errChan := starRequestGen.DoHttpRequest("GetMarathonUrl", nil, nil, nil, "")

	// 		for {
	// 			select {
	// 			// 并发执行http请求，只要name属性和policy.Marathons.MarathonName则退出
	// 			case f := <-finMesos:
	// 				// 保证其他goroutine正常关闭
	// 				finMesos <- f
	// 				// fmt.Println("从ip", v, "获取marathon地址的协程即将完成")
	// 				return
	// 			case err := <-errChan:
	// 				fmt.Println("从messos获取marathon的url失败,错误信息是： ", err, "messos的地址是：", v)
	// 				return
	// 			case response := <-responseMesosChan:
	// 				// 处理response请求
	// 				var messos httpc.Mesoses
	// 				if response == nil {
	// 					break
	// 				}
	// 				byt, err := ioutil.ReadAll(response.Body)
	// 				if err != nil {
	// 					errChan <- err
	// 				}
	// 				defer response.Body.Close()
	// 				err = json.Unmarshal(byt, &messos)
	// 				if err != nil {
	// 					errChan <- err
	// 				}
	// 				if len(messos.Frameworks) == 0 {
	// 					err := errors.New("从messos获取marathon的地址失败，frameworks为空")
	// 					errChan <- err
	// 				}
	// 				for _, framwork := range messos.Frameworks {
	// 					if strings.EqualFold(framwork.Name, policy.Marathons.MarathonName) {

	// 						// 只会执行一次
	// 						once.Do(func() {
	// 							marathonUrl <- framwork.MarathonUrl

	// 						})

	// 						fmt.Println("获取marathon的host成功：", framwork.MarathonUrl)
	// 						finMesos <- true
	// 						break
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}(v, marathonUrl)
	// }

	/***************获取到marathon地址后，执行扩容操作，整个操作设置为发布订阅模式***************/

	go func(rule db.AppScaleRule, finScaled chan bool) {

		errChan := make(chan error, 1)
		responseChan := make(chan *http.Response, 1)
		// 定义扩容的数量
		var scaleInstance int
		for {
			select {
			//已经获取到maramthon的地址
			case m := <-marathonEndpoint:
				// 获取当前的实例数
				fmt.Println("=======marathon地址======", m)
				starRequestGen := httpc.NewStarRequestGen(m, routers.MarathonRouter)
				responseGetAppsChan, errGetAppsChan := starRequestGen.DoHttpRequest("GetApps", httpc.Mapstring{"app_id": rule.AppId}, nil, nil, "")
				errChan = errGetAppsChan
				var marathonApp MarathonApp
				responseGetAppsByt, err := ioutil.ReadAll((<-responseGetAppsChan).Body)
				errChan <- err
				json.Unmarshal(responseGetAppsByt, &marathonApp)
				// 通过marathonUrl 查询当前的实例数目
				nowInstance := marathonApp.Apps.Instances
				// 扩容,只需要更改实例的数量即可
				if strings.EqualFold(rule.ScaleType, "up") {
					// 设置每次扩容的数量
					// per_auto_scale + 当前已经有的实例数 > scale_threshol则target_instances=scale_threshol
					if nowInstance+rule.PerAutoScale > rule.ScaleThreshold {
						scaleInstance = rule.ScaleThreshold
					} else {
						scaleInstance = nowInstance + rule.PerAutoScale
					}
					// 缩容，只需要更改实例的数量即可
				} else if strings.EqualFold(rule.ScaleType, "down") {
					if nowInstance-rule.PerAutoScale < rule.ScaleThreshold {
						scaleInstance = rule.ScaleThreshold
					} else {
						scaleInstance = nowInstance - rule.PerAutoScale
					}
				}
				scale := new(Scale)
				scale.Instance = scaleInstance

				scaleJson, _ := json.Marshal(scale)
				starRequestGen = httpc.NewStarRequestGen(m, routers.MarathonRouter)
				header := map[string]string{"Content-Type": "application/json;charset=utf-8"}
				responseMesosChan, errMessosChan := starRequestGen.DoHttpRequest("UpdateApp", httpc.Mapstring{"app_id": rule.AppId}, bytes.NewBuffer(scaleJson), header, "")
				errChan = errMessosChan
				responseChan = responseMesosChan
			// 发送错误，退出整个goroutine
			case err := <-errChan:
				fmt.Println(rule.AppId, "扩缩容发生错误，错误信息：", err)
				return
			// 发送成功后，退出整个goroutine
			case response := <-responseChan:
				if response.StatusCode == 200 {
					byt, _ := ioutil.ReadAll(response.Body)
					fmt.Println(rule.AppId, "扩缩容成功,扩容数目为", scaleInstance, "扩缩容类型为：", rule.ScaleType, "返回结果是：", string(byt))
					finScaled <- true
					event := utils.StringJoin("扩缩容执行成功,扩缩容类型为: ", rule.ScaleType)
					InsertLog(event, rule.ContinuePeriod, rule)
					return
				} else {

					fmt.Println(rule.AppId, "扩缩容失败", "扩缩容类型为：", rule.ScaleType, "返回错误码：", response.StatusCode)
					event := utils.StringJoin("扩缩容执行失败,扩缩容类型为: ", rule.ScaleType)
					InsertLog(event, rule.ContinuePeriod, rule)
					finScaled <- false
					return
				}
			}
		}
	}(rule, finScaled)
}

func matchJobs(rule db.AppScaleRule, metrics map[string]*httpc.Cmth, counter map[string]chan bool) {

	app := utils.StringJoin(rule.MarathonName, rule.AppId)

	if rule.Memory == 1 {
		flag, err := matchJob(metrics, app, "memory", rule.ScaleType)
		if flag {
			counter["memory"] <- true
			//计数器加一
			fmt.Println("memory信息满足,扩缩容类型为: ", rule.ScaleType)
			event := utils.StringJoin("memory信息满足,扩缩容类型为: ", rule.ScaleType)
			InsertLog(event, len(counter["memory"]), rule)
		} else {
			if err == nil {
				// 计数器清零
				c := len(counter["memory"])
				for i := 0; i < c; i++ {
					<-counter["memory"]
				}
			}
		}
	}

	if rule.Cpu == 1 {

		flag, err := matchJob(metrics, app, "cpu", rule.ScaleType)
		if flag {
			//计数器加一
			counter["cpu"] <- true
			fmt.Println("cpu信息满足,扩缩容类型为: ", rule.ScaleType)
			event := utils.StringJoin("cpu信息满足,扩缩容类型为: ", rule.ScaleType)
			InsertLog(event, len(counter["cpu"]), rule)
		} else {
			if err == nil {
				// 计数器清零
				c := len(counter["cpu"])
				for i := 0; i < c; i++ {
					<-counter["cpu"]
				}
			}
		}
	}

	if rule.Thread == 1 {
		flag, err := matchJob(metrics, app, "thread", rule.ScaleType)
		if flag {
			//计数器加一
			counter["thread"] <- true
			fmt.Println("thread信息满足,扩缩容类型为: ", rule.ScaleType)
			event := utils.StringJoin("thread信息满足,扩缩容类型为: ", rule.ScaleType)
			InsertLog(event, len(counter["thread"]), rule)
		} else {
			if err == nil {
				// 计数器清零
				c := len(counter["thread"])
				for i := 0; i < c; i++ {
					<-counter["thread"]
				}
			}
		}
	}

	if rule.HaQueue == 1 {
		flag, err := matchJob(metrics, app, "ha_queue", rule.ScaleType)
		if flag {
			//计数器加一
			counter["ha_queue"] <- true
			fmt.Println("ha_queue信息满足,扩缩容类型为: ", rule.ScaleType)
			event := utils.StringJoin("ha_queue信息满足,扩缩容类型为: ", rule.ScaleType)
			InsertLog(event, len(counter["ha_queue"]), rule)
		} else {
			if err == nil {
				// 计数器清零
				c := len(counter["ha_queue"])
				for i := 0; i < c; i++ {
					<-counter["ha_queue"]
				}
			}
		}
	}

}

func matchJob(metrics map[string]*httpc.Cmth, app, mtyp, styp string) (bool, error) {

	metricsData := metrics[mtyp]

	if metricsData == nil {
		fmt.Println("=====metricsData由于某种原因变成空指针====")
		return false, nil
	}
	if len(metricsData.Datas.Results) == 0 {
		err := errors.New("获取指标信息失败")
		fmt.Println("从Prometheus查询到的", mtyp, "数据为空,应用信息：", app)
		return false, err
	}
	for _, value := range metricsData.Datas.Results {
		// 转换成float64
		valeOfFloat64, err := strconv.ParseFloat(value.Values[1].(string), 64)

		if err != nil {
			fmt.Println("转成float时出现错误", err)
			return false, err
		}
		// 获取配额信息
		apps := utils.StringJoin(app, mtyp)
		quota, ok := QuotaInfos.Get(apps)
		fmt.Println("配额信息：", quota, "应用信息：", apps)
		if !ok {
			fmt.Println("配额表quato_info为空")
			return false, nil
		}

		if strings.EqualFold(styp, "up") {

			if valeOfFloat64 > quota.(db.QuotaInfo).MaxThreshold*0.01 {
				return true, nil
			}
		} else if strings.EqualFold(styp, "down") {

			if valeOfFloat64 < quota.(db.QuotaInfo).MinThreshold*0.01 {
				return true, nil
			}
		}
	}
	err := errors.New("扩缩容失败")
	return false, err
}
