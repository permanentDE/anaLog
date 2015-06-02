package persistence

import (
	"encoding/json"
	"fmt"
	//"strconv"

	"github.com/influxdb/influxdb/client"

	"go.permanent.de/anaLog/v1/anaLog/logpoint"
	"go.permanent.de/anaLog/v1/anaLog/persistence/influx"
)

func StorePoint(lp logpoint.LogPoint) error {
	return StorePoints(lp)
}

func StorePoints(lps ...logpoint.LogPoint) error {
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
			Name: "logentries_" + lp.Task,
			Time: lp.Time,
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

	return influx.InsertBatch(points)
}

func GetRecurring() (map[string]map[string]map[string]logpoint.LogPoint, error) {
	//------------------task-----identifier----state------------------
	table := make(map[string]map[string]map[string]logpoint.LogPoint)

	results, err := influx.QueryDB("SELECT * FROM /^logentries_.*/;")
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
