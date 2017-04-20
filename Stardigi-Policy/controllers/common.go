package controllers

import (
	"Stardigi-Policy/model/db"
	"Stardigi-Policy/utils"
	"strconv"
	// "sync"
)

type Scale struct {
	Instance int `json:"instances"`
}

type MarathonApp struct {
	Apps App `json:"app"`
}

type App struct {
	Id        string `json:"id"`
	Instances int    `json:"instances"`
}

func InsertLog(event string, count int, rule db.AppScaleRule) {

	s := db.NewSql("mysql")
	countStatus := utils.StringJoin(strconv.Itoa(count), "/", strconv.Itoa(rule.ContinuePeriod))
	s.Insert(&db.ScaleLog{
		MarathonName: rule.MarathonName,
		AppId:        rule.AppId,
		Time:         utils.CurrentimeString(),
		CountStatus:  countStatus,
		Event:        event,
	})
}

func SetChanNull(counter map[string]chan bool) {

	for _, v := range counter {
		l := len(v)
		for i := 0; i < l; i++ {
			<-v
		}
	}
}
