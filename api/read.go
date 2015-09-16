package api

import (
	"time"

	"go.permanent.de/anaLog/logpoint"
	"go.permanent.de/anaLog/persistence"
)

func Find(task, host, state, rawRegex string, timeRangeGTE, timeRangeLTE time.Time, n uint) ([]logpoint.LogPoint, error) {
	return persistence.Find(task, host, state, rawRegex, timeRangeGTE, timeRangeLTE, n)
}
