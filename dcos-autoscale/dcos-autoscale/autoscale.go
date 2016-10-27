package main

import (
	"fmt"
	"github.com/dcos/dcos-autoscale/common"
	"golang.org/x/net/context"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
    "log"
)

func auto_query_list(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	var applist App_list
	var app AppList  
	var appinfo Applist
	result,_ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(result),&applist)
	CheckErr(err)
	app_id_results := SerachAppid(Sql_app_id,applist.MarathonName)
	fmt.Println("app_id_results1111111",app_id_results)
	for k,result :=range app_id_results{
		if app_id_results != nil{
			app.AppId = app_id_results[k]
			watch_retsult  := SerachAppid(Sql_status,applist.MarathonName,result)
			if watch_retsult != nil{
				if (watch_retsult[0] == "1") || (watch_retsult[1] == "1"){
					app.Status = "1"
				}else{
					app.Status = "0"
				}
			}
			count_event_retsult := SerachInfo(Sql_event_status,applist.MarathonName,result)
			if count_event_retsult != nil{
				app.CountStatus = count_event_retsult[0]
				app.EeventDescription = count_event_retsult[1]
			}
		}
		appinfo.APP = append(appinfo.APP,app)
	}
	ret, err := json.Marshal(appinfo)
	CheckErr(err)
	w.Write(ret)
	return nil
}


func auto_query_details(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	var appinfo App_lnfo
	var app     QueryDetail
	result,_ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(result),&appinfo)
	CheckErr(err)
	app_up_results := SerachInfo(Sql_query_info,appinfo.MarathonName,appinfo.Appid)
	if app_up_results != nil{
		app.AppId = appinfo.Appid
		app.MarathonName = appinfo.MarathonName
		app.ScaleUp.CollectPeriod,_= strconv.Atoi(app_up_results[11])
		app.ScaleUp.ContinuePeriod,_ = strconv.Atoi(app_up_results[12])
		app.ScaleUp.CoolDownPeriod,_ = strconv.Atoi(app_up_results[10])
		app.ScaleUp.Switch,_ = strconv.Atoi(app_up_results[9])
		app.ScaleUp.ScaleAmount,_ = strconv.Atoi(app_up_results[4])
		app.ScaleUp.MaxScaleAmount,_ = strconv.Atoi(app_up_results[3])
	}
	app_down_results := SerachInfo(Sql_query_info_down,appinfo.MarathonName,appinfo.Appid)
	if app_down_results != nil{
		app.MarathonName = appinfo.MarathonName
		app.ScaleDown.CollectPeriod,_= strconv.Atoi(app_down_results[11])
		app.ScaleDown.ContinuePeriod,_ = strconv.Atoi(app_down_results[12])
		app.ScaleDown.CoolDownPeriod,_ = strconv.Atoi(app_down_results[10])
		app.ScaleDown.Switch,_ = strconv.Atoi(app_down_results[9])
		app.ScaleDown.ScaleAmount,_ = strconv.Atoi(app_down_results[4])
		app.ScaleDown.MaxScaleAmount,_ = strconv.Atoi(app_down_results[3])
	}
    app_mem_results := SerachInfo(Sql_quota_memory,appinfo.MarathonName,appinfo.Appid)
    if app_mem_results != nil{
    	app.ScaleUp.Mem,_ = strconv.ParseFloat(app_mem_results[0],32)
    	app.ScaleDown.Mem,_ = strconv.ParseFloat(app_mem_results[1],32)
    //	mem,_ := strconv.ParseFloat(app_mem_results[1],32)
	//    app.ScaleDown.Mem = strconv.FormatFloat(app.ScaleDown.Mem,'f',1,64)
    	
    }

    app_cpu_results := SerachInfo(Sql_quota_cpu,appinfo.MarathonName,appinfo.Appid)
    if app_cpu_results != nil{
    	app.ScaleUp.Cpu,_ = strconv.ParseFloat(app_cpu_results[0],32)
    	app.ScaleDown.Cpu,_ = strconv.ParseFloat(app_cpu_results[1],32)
    }

    app_ha_results := SerachInfo(Sql_quota_ha,appinfo.MarathonName,appinfo.Appid)
    if app_ha_results != nil{
    	app.ScaleUp.RequestQueue,_ = strconv.ParseFloat(app_ha_results[0],32)
    	app.ScaleDown.RequestQueue,_ = strconv.ParseFloat(app_ha_results[1],32)
    }

    app_thread_results := SerachInfo(Sql_quota_thread,appinfo.MarathonName,appinfo.Appid)
    if app_thread_results != nil{
    	app.ScaleUp.Thread,_ = strconv.ParseFloat(app_thread_results[0],32)
    	app.ScaleDown.Thread,_ = strconv.ParseFloat(app_thread_results[1],32)
    }

    ret, err := json.Marshal(app)	
	CheckErr(err)
	w.Write(ret)
	return nil
}


func elastic_expansion_configuration(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {
	
	var create CreateExpanionCap
	var app    SetExpanionCap
	result,_ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(result),&create)
	CheckErr(err)

	updata :=  ChangeToInt(create.ScaleUp.Mem,create.ScaleUp.Cpu,create.ScaleUp.Thread,create.ScaleUp.RequestQueue)
	app_info_result_up := ModelinterFace(Sql_scalerule_insert_up,create.MarathonName,create.AppId,"up",create.ScaleUp.ScaleAmount,create.ScaleUp.MaxScaleAmount,updata[0],
		updata[1],updata[2],updata[3],create.ScaleUp.Switch,create.ScaleUp.CoolDownPeriod,create.ScaleUp.CollectPeriod,create.ScaleUp.ContinuePeriod,
	)
	log.Printf("app_info_result_up: %s\n", app_info_result_up)
	
	downdata :=  ChangeToInt(create.ScaleUp.Mem,create.ScaleDown.Cpu,create.ScaleDown.Thread,create.ScaleDown.RequestQueue)
	app_info_result_down := ModelinterFace(Sql_scalerule_insert_down,create.MarathonName,create.AppId,"down",create.ScaleDown.ScaleAmount,create.ScaleDown.MaxScaleAmount,downdata[0],
		downdata[1],downdata[2],downdata[3],create.ScaleDown.Switch,create.ScaleDown.CoolDownPeriod,create.ScaleDown.CollectPeriod,create.ScaleDown.ContinuePeriod,
	)
	log.Printf("app_info_result_down: %s\n", app_info_result_down)

	app_info_result_mem := ModelinterFace(Sql_quota_insert_mem,create.MarathonName,create.AppId,"memory",create.ScaleUp.Mem,create.ScaleDown.Mem)
	log.Printf("app_info_result_mem: %s\n", app_info_result_mem)

	app_info_result_cpu := ModelinterFace(Sql_quota_insert_cpu,create.MarathonName,create.AppId,"cpu",create.ScaleUp.Cpu,create.ScaleDown.Cpu)
	log.Printf("app_info_result_cpu: %s\n", app_info_result_cpu)
	
	app_info_result_thread := ModelinterFace(Sql_quota_insert_thread,create.MarathonName,create.AppId,"thread",create.ScaleUp.Thread,create.ScaleDown.Thread)
	log.Printf("app_info_result_thread: %s\n", app_info_result_thread)
	
	app_info_result_ha := ModelinterFace(Sql_quota_insert_ha,create.MarathonName,create.AppId,"ha_queue",create.ScaleUp.RequestQueue,create.ScaleDown.RequestQueue)
	log.Printf("app_info_result_ha: %s\n", app_info_result_ha)


	app.Status = "create successfu"
	app.Msgup = "up ok"
	app.Msgdown = "down ok"
	ret, err := json.Marshal(app)	
	CheckErr(err)
	w.Write(ret)
	return nil
}


func updating_elastic_expansion(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	var update UpdateExpaionCap
	var app    UpExpanionCap
	result,_ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(result),&update)
	CheckErr(err)

	updata :=  ChangeToInt(update.ScaleUp.Mem,update.ScaleUp.Cpu,update.ScaleUp.Thread,update.ScaleUp.RequestQueue)
	app_info_result_up := ModelinterFace(Sql_scalerule_update_up,update.MarathonName,update.AppId,"up",update.ScaleUp.ScaleAmount,update.ScaleUp.MaxScaleAmount,
		updata[0],updata[1],updata[2],updata[3],update.ScaleUp.Switch,update.ScaleUp.CoolDownPeriod,update.ScaleUp.CollectPeriod,update.ScaleUp.ContinuePeriod,
	update.MarathonName,update.AppId)
    log.Printf("app_info_result_up: %s\n", app_info_result_up)

	downdata :=  ChangeToInt(update.ScaleUp.Mem,update.ScaleDown.Cpu,update.ScaleDown.Thread,update.ScaleDown.RequestQueue)
	app_info_result_down := ModelinterFace(Sql_scalerule_update_down,update.MarathonName,update.AppId,"down",update.ScaleDown.ScaleAmount,update.ScaleDown.MaxScaleAmount,
		downdata[0],downdata[1],downdata[2],downdata[3],update.ScaleDown.Switch,update.ScaleDown.CoolDownPeriod,update.ScaleDown.CollectPeriod,update.ScaleDown.ContinuePeriod,
	update.MarathonName,update.AppId)
	log.Printf("app_info_result_down: %s\n", app_info_result_down)
	
	app_info_result_mem := ModelinterFace(Sql_quota_update_mem,update.MarathonName,update.AppId,"memory",update.ScaleUp.Mem,update.ScaleDown.Mem)
    log.Printf("app_info_result_mem: %s\n", app_info_result_mem)

	app_info_result_cpu := ModelinterFace(Sql_quota_update_cpu,update.MarathonName,update.AppId,"cpu",update.ScaleUp.Cpu,update.ScaleDown.Cpu)
    log.Printf("app_info_result_cpu: %s\n", app_info_result_cpu)

	app_info_result_thread := ModelinterFace(Sql_quota_update_thread,update.MarathonName,update.AppId,"thread",update.ScaleUp.Thread,update.ScaleDown.Thread)
	log.Printf("app_info_result_thread: %s\n", app_info_result_thread)

	app_info_result_ha := ModelinterFace(Sql_quota_update_ha,update.MarathonName,update.AppId,"ha_queue",update.ScaleUp.RequestQueue,update.ScaleDown.RequestQueue)
	log.Printf("app_info_result_ha: %s\n", app_info_result_ha)

	app.Status = "update successfu"
	app.Msgup = "up ok"
	app.Msgdown = "down ok"
	ret, err := json.Marshal(app)	
	CheckErr(err)
	w.Write(ret)
	return nil
}

func suspend(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	var appsuspend Suspend
	var app SuspendInfo
	result,_ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(result),&appsuspend)
	CheckErr(err)
	suspend_retsult_up := ModelinterFace(Sql_update_watch,appsuspend.MarathonName,appsuspend.AppId)
	if strconv.FormatInt(suspend_retsult_up,10) != "" {
		app.Status = "update ok"
		app.Msgup = "up rule pause successfu"
	}


	suspend_retsult_down := ModelinterFace(Sql_update_watch_down,appsuspend.MarathonName,appsuspend.AppId)
	if strconv.FormatInt(suspend_retsult_down,10) != "" {
		app.Status = "update ok"
		app.Msgdown = "down rule pause successful"
	}

	ret, err := json.Marshal(app)	
	CheckErr(err)
	w.Write(ret)
	return nil

}

func turn_on(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	var appstart Start
	var app StartInfo
	result,_ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(result),&appstart)
	CheckErr(err)
	start_retsult_up := ModelinterFace(Sql_update_start,appstart.MarathonName,appstart.AppId)
	if strconv.FormatInt(start_retsult_up,10) != "" {
		app.Status = "update ok"
		app.Msgup = "up rule recover successfu"
	}

	start_retsult_down := ModelinterFace(Sql_start_update_watch,appstart.MarathonName,appstart.AppId)
	if strconv.FormatInt(start_retsult_down,10) != "" {
		app.Status = "update ok"
		app.Msgdown = "down rule recover successful"
	}

	ret, err := json.Marshal(app)	
	CheckErr(err)
	w.Write(ret)
	return nil

}

