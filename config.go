package main

import (
	"fmt"
	"go.uber.org/config"
)

type BiliConfig struct {
	Cookie string `yaml:"cookie"`
}

type AppConfig struct {
	Debug  bool   `yaml:"debug"`
	Output string `yaml:"output"`
}

type Config struct {
	BiliConfig *BiliConfig `yaml:"biliconfig"`
	AppConfig  *AppConfig  `yaml:"appconfig"`
}

func GetConfig() *Config {
	c, err := config.NewYAML(config.File("config.yaml"))
	if err != nil {
		fmt.Println(err.Error())
	}

	var cfg *Config
	err = c.Get("").Populate(&cfg)
	if err != nil {
		fmt.Println(err.Error())
	}

	return cfg
}
