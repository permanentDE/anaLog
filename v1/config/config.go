package config

import (
	"code.google.com/p/gcfg"
	"flag"
	"go.iondynamics.net/iDlogger"
)

var Std *Config

type Config struct {
	AnaLog struct {
		Listen         string
		Fcgi           bool
		Worker         uint8
		Workspace      string
		CookieSecret   string
		DevelopmentEnv bool
	}

	Influx struct {
		Host     string
		Port     uint16
		Database string
		User     string
		Pass     string
	}

	AppLog struct {
		SlackLogUrl string
	}
}

func init() {
	ConfigPath := flag.String("conf", "./config.gcfg", "Specify a file containing the configuration or \"./config.gcfg\" is used")
	flag.Parse()
	Std = &Config{}
	err := gcfg.ReadFileInto(Std, *ConfigPath)

	if err != nil {
		iDlogger.Emerg("config: ", err)
	}
}
