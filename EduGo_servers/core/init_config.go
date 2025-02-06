package core

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
)

var confPath = "config.yaml"

type System struct {
	Ip string `yaml:"ip"`
	Port int `yaml:"port"`
}
type Config struct {
	System System `yaml:"system"`
}

func ReadConfig() {
	byteDate, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(byteDate, &config)
	if err != nil {
		panic(fmt.Sprintf("setting文件yaml格式错误: %v", err))
	}
	fmt.Println(config)
}