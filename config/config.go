package config

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/config"
	"log"
)

var C = DefaultConfig()

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

func DefaultConfig() *Config {
	return &Config{
		BiliConfig: &BiliConfig{Cookie: ""},
		AppConfig: &AppConfig{
			Debug:  false,
			Output: "./",
		},
	}
}

func GetConfig(path cli.Path) *Config {
	c, err := config.NewYAML(config.File(path))
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg := &Config{
		BiliConfig: &BiliConfig{Cookie: ""},
		AppConfig: &AppConfig{
			Debug:  false,
			Output: "./",
		},
	}
	err = c.Get("").Populate(&cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	C = cfg
	return cfg
}
