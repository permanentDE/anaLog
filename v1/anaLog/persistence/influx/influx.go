package influx

import (
	"fmt"
	influx "github.com/influxdb/influxdb/client"
	idl "go.iondynamics.net/iDlogger"
	"net/url"

	"go.permanent.de/anaLog/v1/config"
)

var db string

func connect2influx() *influx.Client {
	db = config.Influx.Database

	u, err := url.Parse(fmt.Sprintf("http://%s:%d", config.Influx.Host, config.Influx.Port))
	if err != nil {
		idl.Emerg(err)
	}

	conf := influx.Config{
		URL:      *u,
		Username: config.Influx.User,
		Password: config.Influx.Pass,
	}

	con, err := influx.NewClient(conf)
	if err != nil {
		idl.Emerg(err)
	}

	return con
}
