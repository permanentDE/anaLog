package influx

import (
	"encoding/json"
	"fmt"

	"github.com/influxdb/influxdb/client"

	"go.permanent.de/anaLog/v1/anaLog/logpoint"
)

type Adapter struct {
	client *client.Client
}

func GetAdapter() *Adapter {
	return &Adapter{client: connect2influx()}
}

func (a *Adapter) Close() {
	//influx http api client doesn't have to be closed
}

func (a *Adapter) queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, rerr := a.client.Query(q); err == nil {
		if rerr != nil {
			err = rerr
			return
		}
		if response.Error() != nil {
			err = response.Error()
			return
		}
		res = response.Results
	}
	return
}

func (a *Adapter) insert(p client.Point) error {
	return a.insertBatch([]client.Point{p})
}

func (a *Adapter) insertBatch(ps []client.Point) error {
	bps := client.BatchPoints{
		Points:   ps,
		Database: db,
	}

	_, err := a.client.Write(bps)

	return err
}

func (a *Adapter) StorePoint(lp logpoint.LogPoint) error {
	return a.StorePoints(lp)
}

func (a *Adapter) StorePoints(lps ...logpoint.LogPoint) error {
	var points []client.Point

	for _, lp := range lps {

		byt, err := json.Marshal(lp)
		if err != nil {
			return err
		}

		fields := map[string]interface{}{
			"task":       lp.Task,
			"identifier": lp.RunId,
			"logpoint":   string(byt),
		}

		if lp.Message != "" {
			fields["message"] = lp.Message
		}

		if lp.Raw != "" {
			fields["raw"] = lp.Raw
		}

		if len(lp.Data) > 0 {
			fields["data"] = lp.Data
		}

		points = append(points, client.Point{
			Measurement: "logentries_" + lp.Task,
			Time:        lp.Time,
			Tags: map[string]string{
				"host":     lp.Host,
				"mode":     fmt.Sprint(lp.Mode),
				"priority": fmt.Sprint(lp.Priority),
				"state":    fmt.Sprint(lp.State),
				//"analyzed": strconv.FormatBool(lp.Analyzed),
			},
			Fields:    fields,
			Precision: "n",
		})
	}

	return a.insertBatch(points)
}

func (a *Adapter) GetRecurring() (map[string]map[string]map[string]logpoint.LogPoint, error) {
	table := make(map[string]map[string]map[string]logpoint.LogPoint)
	//----------------task-----identifier----state------------------

	results, err := a.queryDB("SELECT * FROM /^logentries_.*/;")
	if err != nil {
		return table, err
	}

	for _, result := range results {
		for _, row := range result.Series {
			columns := make(map[string]int)
			for key, column := range row.Columns {
				columns[column] = key
			}

			for _, fields := range row.Values {
				lpbytes := []byte(fields[columns["logpoint"]].(string))
				var lp logpoint.LogPoint
				err = json.Unmarshal(lpbytes, &lp)
				if err != nil {
					return table, err
				}

				_, ok := table[lp.Task]
				if ok {
					_, ok = table[lp.Task][lp.RunId]
					if !ok {
						table[lp.Task][lp.RunId] = make(map[string]logpoint.LogPoint)
					}
				} else {
					table[lp.Task] = make(map[string]map[string]logpoint.LogPoint)
					table[lp.Task][lp.RunId] = make(map[string]logpoint.LogPoint)
				}

				table[lp.Task][lp.RunId][fmt.Sprint(lp.State)] = lp

			}
		}
	}

	return table, err
}

func (a *Adapter) StoreAnalysisRCserial(lrc []byte) error {
	return fmt.Errorf("%s", "not implemented")
}

func (a *Adapter) GetLatestAnalysisRCserial() ([]byte, error) {
	return nil, fmt.Errorf("%s", "not implemented")
}

func (a *Adapter) GetEndByStart(begin logpoint.LogPoint) (logpoint.LogPoint, error) {
	return logpoint.LogPoint{}, fmt.Errorf("%s", "not implemented")
}

func (a *Adapter) GetLastBegin(taskname string) (logpoint.LogPoint, error) {
	return logpoint.LogPoint{}, fmt.Errorf("%s", "not implemented")
}
