package influx

import (
	"fmt"
	influx "github.com/influxdb/influxdb/client"
	idl "go.iondynamics.net/iDlogger"
	"net/url"

	"go.permanent.de/anaLog/v1/config"
)

var influxClient *influx.Client

func connect2influx() {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", config.Std.Influx.Host, config.Std.Influx.Port))
	if err != nil {
		idl.Emerg(err)
	}

	conf := influx.Config{
		URL:      *u,
		Username: config.Std.Influx.User,
		Password: config.Std.Influx.Pass,
	}

	con, err := influx.NewClient(conf)
	if err != nil {
		idl.Emerg(err)
	}

	influxClient = con
}

func getInflux() *influx.Client {
	if influxClient == nil {
		connect2influx()
	}
	return influxClient
}

func QueryDB(cmd string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: config.Std.Influx.Database,
	}
	if response, err := getInflux().Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

func Insert(p influx.Point) error {
	return InsertBatch([]influx.Point{p})
}

func InsertBatch(ps []influx.Point) error {
	bps := influx.BatchPoints{
		Points:   ps,
		Database: config.Std.Influx.Database,
	}

	_, err := getInflux().Write(bps)

	return err
}
