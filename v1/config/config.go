package config

import (
	"flag"

	"code.google.com/p/gcfg"
	"github.com/badgerodon/socketmaster/protocol"
	"go.iondynamics.net/iDlogger"
)

func init() {
	ConfigPath := flag.String("conf", "./config.gcfg", "Specify a file containing the configuration or \"./config.gcfg\" is used")
	flag.Parse()
	cfg := &Config{}
	err := gcfg.ReadFileInto(cfg, *ConfigPath)

	if err != nil {
		iDlogger.Emerg("config: ", err)
	}

	if cfg.AnaLogCfg.UseSocketMaster && cfg.AnaLogCfg.Fcgi {
		iDlogger.Emerg("Invalid configuration: Can't serve via FastCGI AND SocketMaster simultaneously. Disable one of them!")
	}

	if cfg.AnaLogCfg.UseSocketMaster && cfg.AnaLogCfg.Listen != "" {
		iDlogger.Emerg("Invalid configuration: Can't serve via inbuilt HTTP Server AND SocketMaster simultaneously. Disable one of them!")
	}

	if !(cfg.AnaLogCfg.UseInflux || cfg.AnaLogCfg.UseMongo) || (cfg.AnaLogCfg.UseInflux && cfg.AnaLogCfg.UseMongo) {
		iDlogger.Emerg("Invalid configuration: Exactly one database backend has to be defined")
	}

	fillVars(cfg)
}

func fillVars(cfg *Config) {
	AnaLog = cfg.AnaLogCfg
	Influx = cfg.InfluxCfg
	Mongo = cfg.MongoCfg
	SocketMaster = cfg.SocketMasterCfg
	AppLog = cfg.AppLogCfg
}

type Config struct {
	AnaLogCfg       `gcfg:"AnaLog"`
	InfluxCfg       `gcfg:"Influx"`
	MongoCfg        `gcfg:"Mongo"`
	SocketMasterCfg `gcfg:"SocketMaster"`
	AppLogCfg       `gcfg:"AppLog"`
}

type AnaLogCfg struct {
	UseInflux         bool
	UseMongo          bool
	SchedulerInterval string
	Listen            string
	Fcgi              bool
	UseSocketMaster   bool
	Worker            uint8
	Workspace         string
	CookieSecret      string
	NagiosSecret      string
	Domain            string
	DevelopmentEnv    bool
}

var AnaLog AnaLogCfg

type InfluxCfg struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Pass     string
}

var Influx InfluxCfg

type MongoCfg struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Pass     string
}

var Mongo MongoCfg

type SocketMasterCfg struct {
	protocol.SocketHTTPDefinition
	protocol.SocketTLSDefinition

	Port int
}

var SocketMaster SocketMasterCfg

type AppLogCfg struct {
	SlackLogUrl string
}

var AppLog AppLogCfg
