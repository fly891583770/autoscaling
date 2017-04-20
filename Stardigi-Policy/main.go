package main

import (
	"Stardigi-Policy/config"
	"Stardigi-Policy/controllers"
	"Stardigi-Policy/model/db"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

var pConfig *string = flag.String("pConfig", "conf/stardigi_policy.json", "configuration of stardigi-policy")
var log *string = flag.String("logfile", "/home/paas/autoscaling.log", "log file")

func main() {
	flag.Parse()
	policy, err := config.LoadPolicyConfig(*pConfig)
	if err != nil {
		fmt.Println("===err==", err)
	}
	db.NewDBClient(policy.Db.Type, policy.Db.Datasource)
	go controllers.PRun()
	go controllers.MRun(policy)
	go controllers.SRun(policy)
	Run()
}

func Run() {
	exit := make(chan os.Signal, 0)
	signal.Notify(exit, os.Kill, os.Interrupt)
	for {
		select {
		case <-exit:
			fmt.Println("main函数退出")
			return
		}
	}
}
