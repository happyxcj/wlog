package main

import (
	"github.com/happyxcj/wlog"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"fmt"
)

func ParseLogConfig(confFilePath string) (wlog.EasyMultiConfig, error) {
	var conf wlog.EasyMultiConfig
	content, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(content, &conf)
	return conf, err
}

func main() {
	cfg, err := ParseLogConfig("log_config.yaml")
	if err != nil {
		panic(fmt.Sprint("failed to parse config file, error: ", err.Error()))
	}
	//todo
	fmt.Println("======",cfg.MultiPaths)
	loggers, err := cfg.Create()
	if err != nil {
		panic(fmt.Sprint("failed to create logger, error: ", err.Error()))
	}
	defer func() {
		for _, logger := range loggers {
			logger.Close()
		}
	}()
	var logger  *wlog.Logger
	for i := 0; i < 30; i++ {
		if i<10{
			logger=loggers[0]
		}else if i<20{
			logger=loggers[1]
		}else{
			logger=loggers[2]
		}
		logger.With(wlog.String("username", "root"),
			wlog.String("password", "root")).
			Info("failed to login account ", i)
	}
}
