package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
//	"net/http"
)

func ReadJson(filename string) (map[string]string, error) {
	var config = map[string]string{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return config, nil
}



func LoadMacConfig (fileName string) (*MacConfig, error) {

	var mac MacConfig
	byt , err := ioutil.ReadFile(fileName)
//	fmt.Printf("%s",byt)
	if err != nil {
		return &mac, err
		fmt.Println(err)
	}
	
	if err := json.Unmarshal(byt, &mac); err != nil {
		// 输出错误日志
		return &mac, err
	//	fmt.Printf("%s",mac.Datasource)
	}
//	fmt.Printf("%s",mac.Datasource)
	return &mac, nil
}


func ChangeToInt(args ...float32) []int{

	var info = make([]int,0,len(args))
	for _,v :=range args{
		fmt.Println("v:\n",v)

		if v > 0 {
			info = append(info,1)
		}else{
			info = append(info,0)
		}

	}
	fmt.Println(info)
	return info 
}

/*
func CheckErr(err error) []byte{

	var app ErrorStatus
	var ret []byte
	defer ret = func(ret []byte) []byte{
		if err := recover;err != nil{	
			app.Status = "NOT OK"
			ret, err = json.Marshal(app)	
		}
	}(ret)

	if err != nil {
		panic(err)
	}
    return ret
}

*/
func  CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
