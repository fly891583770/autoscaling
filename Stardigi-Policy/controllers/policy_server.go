// policy_server service as a server runs as a background process
// use to get policy rule from db
// the rule transfer in a channel will received by scale_server process
package controllers

import (
	"Stardigi-Policy/model/db"
	"Stardigi-Policy/utils"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

/**
app_scale_rule
根据switch=1(打开了扩缩容)，查询所有需要扩缩容的应用
接着已switch=1为条件查询scale_type判断up->扩容，down->缩容，
对于扩容/缩容，判断cpu/memory/ha_queue/thread是否为1(只要有一个打开就执行扩缩容)
per_auto_scale是每次扩容的实例数，scale_threshold是指当前实例数的最大限制，marathon接口查询

quota_info
根据app_scale_rule查出的需要扩缩容的marathon_name,app_id找到扩缩容的条件
>max_threshod 扩容，
<min_threshold 缩容
*/

var (
	RuleUpChan   = make(chan []db.AppScaleRule, 1)
	RuleDownChan = make(chan []db.AppScaleRule, 1)
	// QuotaInfos   = make(map[string]db.QuotaInfo, 0)
	QuotaInfos *utils.SyncMap
	mutex      sync.Mutex
)


func init() {
	QuotaInfos = utils.NewSyncMap()
}

func PRun() {
	// 收集错误信息，写入日志

	error_c := make(chan error, 1)

	// 启动后台程序判断是否发生错误

	go func() {
		signals := make(chan os.Signal, 0)
		signal.Notify(signals, os.Interrupt, os.Kill)

		// 初始化
		for {
			select {
			case <-time.After(5 * time.Second):
				// finishRule由scale模块控制
				if len(RuleUpChan) == 0 && len(RuleDownChan) == 0 {
					fmt.Println("==开始查询规则===")
					GetRuleFromDB(RuleUpChan, RuleDownChan, error_c)
				}

			case <-signals:
				fmt.Println("由于系统执行了中断命令，规则收集程序退出...")
				return
			case e := <-error_c:
				fmt.Println("规则收集发生错误，错误信息：", e)
			}
		}
	}()
}

func GetRuleFromDB(ruleUpChan, ruleDownChan chan []db.AppScaleRule, error_c chan error) {

	s := db.NewSql("mysql")

	var waitMysql sync.WaitGroup
	// 表示需要查询规则和配额表
	waitMysql.Add(2)

	appScaleRule := new(db.AppScaleRule)
	var appScaleRules []db.AppScaleRule
	var waitRule sync.WaitGroup
	waitRule.Add(2)
	go func(appScaleRule *db.AppScaleRule, appScaleRules []db.AppScaleRule, wait *sync.WaitGroup, error_c chan error) {
		// 由于等待数据库完成查询和分离扩缩容信息
		// switch == 1 means scale enable
		marathon_name := os.Getenv("marathon_name")
		err := s.QueryTable(appScaleRule).Filter("switch", "1").Filter("marathon_name", marathon_name).All(&appScaleRules)
		if err != nil {
			fmt.Println("查询数据库出错", err)
			waitMysql.Done()
			error_c <- err
			return
		}
		// 计算up和down的存在情况
		countUpAndDown := 0
		countUp := 0
		countDown := 0
		var onceUp sync.Once
		var onceDown sync.Once
		addf := func() {
			countUpAndDown++
		}
		for _, r := range appScaleRules {
			if strings.EqualFold(r.ScaleType, "up") {
				onceUp.Do(addf)
				countUp++
			} else if strings.EqualFold(r.ScaleType, "down") {
				onceDown.Do(addf)
				countDown++
			}
		}
		if countUpAndDown == 0 {
			waitRule.Done()
			waitRule.Done()
		} else if countUpAndDown == 1 {
			waitRule.Done()
		}

		// 否则表示查询完成，开始分离扩容和缩容
		if countUp > 0 {
			go selectRules(appScaleRules, "up", ruleUpChan, &waitRule)

		}

		// 分离出缩容记录 ScaleType = down
		if countDown > 0 {
			go selectRules(appScaleRules, "down", ruleDownChan, &waitRule)

		}

	}(appScaleRule, appScaleRules, &waitRule, error_c)

	waitRule.Wait()
	waitMysql.Done()

	go func(error_c chan error) {

		// 查询配额信息
		quotaInfo := new(db.QuotaInfo)
		var quotaInfos []db.QuotaInfo
		// 查询扩容的配额信息
		err := s.QueryTable(quotaInfo).All(&quotaInfos)
		if err != nil {
			waitMysql.Done()
			error_c <- err
			return
		}
		// 格式化配额信息
		formatQuota(QuotaInfos, quotaInfos)
		waitMysql.Done()
	}(error_c)

	// 等待规则和配额查询完成
	waitMysql.Wait()
}

func formatQuota(quota *utils.SyncMap, quotaInfos []db.QuotaInfo) {

	for _, value := range quotaInfos {
		appInfo := utils.StringJoin(value.MarathonName, value.AppId, value.RuleType)
		quota.Set(appInfo, value)
	}
}

func selectRules(appScaleRules []db.AppScaleRule, name string, appScaleRule chan []db.AppScaleRule, wait *sync.WaitGroup) { //container *[]db.AppScaleRule) {

	rule := make([]db.AppScaleRule, 0)
	for _, r := range appScaleRules {
		if strings.EqualFold(r.ScaleType, name) {
			rule = append(rule, r)
		}
	}
	mutex.Lock()
	appScaleRule <- rule
	mutex.Unlock()
	wait.Done()

	return
}
