package config

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/config"
	"log"
)

var C = DefaultConfig()

type Config struct {
	Cookie string `yaml:"cookie"`
	Debug  bool   `yaml:"debug"`
	Output string `yaml:"output,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Cookie: "",
		Debug:  false,
		Output: "",
	}
}

func GetConfig(path cli.Path) *Config {
	c, err := config.NewYAML(config.File(path))
	if err != nil {
		log.Fatal(err.Error())
	}

	var cfg Config
	err = c.Get("").Populate(&cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	C = &cfg
	return &cfg
}
