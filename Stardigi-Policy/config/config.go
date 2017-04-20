package config

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func LoadPolicyConfig(fileName string) (*PolicyConfig, error) {

	var config PolicyConfig
	byt, err := ioutil.ReadFile(fileName)
	fmt.Println("---------byt----",string(byt))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byt, &config)
	fmt.Println("-----------config-------",config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

/**
 * 采用反射获取配置文件
 */
func LoadConfig(container interface{}) {

}
