package controllers

import (
	"Stardigi-Policy/config"
	"Stardigi-Policy/model/db"
	"Stardigi-Policy/model/httpc"
	"Stardigi-Policy/routers"
	"Stardigi-Policy/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	MetricsOld = map[string]string{
		"cpu":      `avg by (marathon_app_id) (irate(container_cpu_system_seconds_total{marathon_app_id='{{app_id}}',id=~'^/docker.*'}[1m])) + avg by (marathon_app_id) (irate(container_cpu_user_seconds_total{marathon_app_id='{{app_id}}',id=~'^/docker.*'}[1m]))`,
		"memory":   `avg by (marathon_app_id)(container_memory_usage_bytes{marathon_app_id='{{app_id}}',id=~'^/docker.*'}/container_spec_memory_limit_bytes{marathon_app_id='{{app_id}}',id=~'^/docker.*'})`,
		"ha_queue": `haproxy_server_current_queue{backend=~'^{{app_id}}_{{marathon_name}}.*'}`,
		"thread":   `tomcat_threadRunning{appId='{{app_id}}'}`,
	}
)

var (
	MetricsUp   = make(chan map[string]map[string]*httpc.Cmth, 0)
	MetricsDown = make(chan map[string]map[string]*httpc.Cmth, 0)
	RulesUp     = make(chan db.AppScaleRule, 0)
	RulesDown   = make(chan db.AppScaleRule, 0)
)

func MRun(policy *config.PolicyConfig) {

	// 接收规则信息，定时判断是否满足条件

	go func() {
		signals := make(chan os.Signal, 0)
		signal.Notify(signals, os.Interrupt, os.Kill)
		errorsMrun := make(chan error, 1)
		for {
			select {
			// 只要scale模块要求查询数据，则执行操作
			case rulesUp := <-RulesUp:
				fmt.Println("==开始从Prometheus查询扩容需要的指标信息==")
				RuleJob(policy, rulesUp, errorsMrun)
			case rulesDown := <-RulesDown:
				fmt.Println("==开始从Prometheus查询缩容需要的指标信息==")
				RuleJob(policy, rulesDown, errorsMrun)
			case err := <-errorsMrun:
				fmt.Println("查询指标信息发生错误，错误信息：", err)
			case <-signals:
				fmt.Println("由于系统执行了中断命令，指标收集程序退出...")
				return
			}
		}
	}()
}

func RuleJob(policy *config.PolicyConfig, rule db.AppScaleRule, errorsMrun chan error) {
	// 设置metrics完整性校验
	var waitRule sync.WaitGroup
	waitRule.Add(1)
	// 每条规则启动协程收集数据
	rescmth := make(map[string]map[string]*httpc.Cmth, 0)

	urlParam := map[string]string{"{{app_id}}": rule.AppId, "{{marathon_name}}": rule.MarathonName}
	metricDest := make(map[string]string, 0)
	SetMetricUrl(MetricsOld, metricDest, urlParam)
	GetMetricFromPrometheus(policy, rule, metricDest, rescmth, &waitRule, errorsMrun)
	waitRule.Wait()
	if strings.EqualFold(rule.ScaleType, "up") {
		MetricsUp <- rescmth
	} else if strings.EqualFold(rule.ScaleType, "down") {
		MetricsDown <- rescmth
	}
}

/**
 * 考虑到后期会有大量的请求，获取的过程采用异步HTTP方式
 */
func GetMetricFromPrometheus(policy *config.PolicyConfig, rule db.AppScaleRule, metricDest map[string]string, rescmth map[string]map[string]*httpc.Cmth, waitRule *sync.WaitGroup, errorsMrun chan error) {
	//创建prometheus的api
	// 只需要运行一次once函数
	host := utils.StringJoin("http://", policy.Prometheuss.PrometheusAddr, ":", policy.Prometheuss.PrometheusPort)
	starRequestGen := httpc.NewStarRequestGen(host, routers.ProRouter)
	// 存储查询结果
	// 同时必须写一个循环判断map内结果是否完整
	metrics := make(map[string]*httpc.Cmth)
	// metrics := new(utils.SyncMap)
	// 为每个查询启动一个协程，并通过通道去控制查询是否结束

	var lock sync.RWMutex
	for k, v := range metricDest {
		fmt.Println("prometheus地址：", host, "类型：", k, "url:", v)
		go func(metrics map[string]*httpc.Cmth, k, v string) {
			responseChan, errChan := starRequestGen.DoHttpRequest("GetMetric", nil, nil, nil, v)
			for {
				select {
				case response := <-responseChan:
					// 容错1：此处容易会概率性出现空指针异常，response为空
					if response == nil {
						err := errors.New("查询metrics出错，response为空")
						errChan <- err
						return
					}
					// json转bean
					body, err := ioutil.ReadAll(response.Body)
					fmt.Println("========metric_server=========", string(body), "=>", k)
					if err != nil {
						errorsMrun <- err
					}
					cmth := new(httpc.Cmth)
					err = json.Unmarshal(body, cmth)
					if err != nil {
						errorsMrun <- err
					}
					// 注意此处传入的指针类型
					// 需保证使用者也是这种类型
					lock.Lock()
					metrics[k] = cmth
					lock.Unlock()
					return
				case err := <-errChan:
					fmt.Println("查询metrics失败，错误信息：", err, "app信息是：", rule.MarathonName, rule.AppId, rule.ScaleType)
					errorsMrun <- err
					return
				case <-time.After(60 * time.Second):
					err := errors.New("get response from http timeout")
					errorsMrun <- err
					return
				}
			}
		}(metrics, k, v)
	}
	//  判断查询结果是否完整
	go func(rule db.AppScaleRule) {
		for {
			select {
			case <-time.After(1 * time.Millisecond):
				if len(metrics) == 4 {
					lock.RLock()
					appInfo := utils.StringJoin(rule.MarathonName, rule.AppId, rule.ScaleType)
					rescmth[appInfo] = metrics
					lock.RUnlock()
					waitRule.Done()
					return
				}
			}
		}
	}(rule)

}

/**
 * 定义一个通道，写入appi_id，
  f_app_id = app_id.replace("/","_")
    f_app_id = f_app_id[1:]
"query=haproxy_server_current_queue%7Bbackend%3D~%27%5E" + f_app_id + "_"+marathon_name+".*%27%7D"
"ha_queue": `haproxy_server_current_queue{backend=~'^{{app_id}}_{{marathon_name}}.*'}`,
urlParam := map[string]string{"{{app_id}}": rule.AppId, "{{marathon_name}}": rule.MarathonName}
*/
func SetMetricUrl(origin map[string]string, dest map[string]string, param map[string]string) {

	for mk, mv := range origin {
		for pk, pv := range param {
			if strings.EqualFold(mk, "ha_queue") {
				pv = strings.Replace(pv, "/", "_", -1)
				pv = pv[1:]
			}
			mv = strings.Replace(mv, pk, pv, -1)
		}
		urlEncode := utils.UrlBase64(mv)
		dest[mk] = urlEncode
	}
}
