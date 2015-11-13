package api

import (
	"fmt"
	"time"

	"go.permanent.de/anaLog/alertBlocker"
	"go.permanent.de/anaLog/analysis"
	"go.permanent.de/anaLog/logpoint"
	"go.permanent.de/anaLog/persistence"
)

type HumanReadableResult struct {
	analysis.Result
	AvgHr            string
	StdDevHr         string
	AvgQrHr          string
	StdDevQrHr       string
	IntervalAvgHr    string
	IntervalStdDevHr string
}

func Find(task, runId, host, state, rawRegex string, timeRangeGTE, timeRangeLTE time.Time, n uint) ([]logpoint.LogPoint, error) {
	return persistence.Find(task, runId, host, state, rawRegex, timeRangeGTE, timeRangeLTE, n)
}

func Results() (hrc map[string]HumanReadableResult, err error) {
	rc := analysis.NewResultContainer()
	err = rc.LoadLatest()
	if err != nil {
		return
	}

	hrc = make(map[string]HumanReadableResult)

	rc.Range(func(task string, res analysis.Result) {
		hrc[task] = HumanReadableResult{
			Result:           res,
			AvgHr:            fmt.Sprint(res.Avg),
			StdDevHr:         fmt.Sprint(res.StdDev),
			AvgQrHr:          fmt.Sprint(res.AvgQr),
			StdDevQrHr:       fmt.Sprint(res.StdDevQr),
			IntervalAvgHr:    fmt.Sprint(res.IntervalAvg),
			IntervalStdDevHr: fmt.Sprint(res.IntervalStdDev),
		}
	})

	return
}

func Problems() map[string][]string {
	return alertBlocker.AllProblems()
}
