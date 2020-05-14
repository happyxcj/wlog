package main

import (
	"github.com/happyxcj/wlog"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"fmt"
)

func ParseLogConfig(confFilePath string) (wlog.Config, error) {
	var conf wlog.Config
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
	logger, err := cfg.Create()
	if err != nil {
		panic(fmt.Sprint("failed to create logger, error: ", err.Error()))
	}
	defer logger.Close()
	for i := 0; i < 500; i++ {
		logger.With(wlog.String("username", "root"),
			wlog.String("password", "root")).
			Info("failed to login account ", i)
	}
}
