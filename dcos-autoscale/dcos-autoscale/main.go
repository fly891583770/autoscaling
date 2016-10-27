package main

import (
	"github.com/codegangsta/cli"
	"github.com/dcos/dcos-autoscale/common"
	"golang.org/x/net/context"
	"log"
	"os"
	"fmt"
	"flag"
)

// Init database before
var macConfig *string = flag.String("macConfig","conf.json","configuration of machine")
// CREATE TABLE images(id varchar(36) NOT NULL,image_id varchar(120) NOT NULL, visibility boolean, owner varchar(120), submission_date DATE, image_name varchar(120) NOT NULL, label varchar(120), PRIMARY KEY(owner,image_name,label));
var macconfig  *MacConfig

func main() {
	serveCommand := cli.Command{
		Name:      "serve",
		ShortName: "s",
		Usage:     "Serve the API",
		Flags:     []cli.Flag{common.FlAddr},
		Action:    action(serveAction),
	}

//	fmt.Println(common.FlAddr.Value)
	log.SetPrefix("==>")
	log.Printf("listen: %v", common.FlAddr.Value)
	flag.Parse()	
	var err error
	macconfig,err = LoadMacConfig(*macConfig)
	if err != nil{
		fmt.Println("err:",err)
		log.Printf("LoadMacConfig is error")
	}

	fmt.Printf("%s\n",macconfig.Datasource)
	common.Run("dcos-autoscale", serveCommand)
}

func serveAction(c *cli.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "addr", c.String("addr"))
	return common.ServeCmd(c, ctx, routes)
}

func action(f func(c *cli.Context) error) func(c *cli.Context) {
	return func(c *cli.Context) {
		err := f(c)
		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
	}
}
