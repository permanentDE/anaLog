package persistence

import (
	"errors"
	"sort"
	"time"

	"go.permanent.de/anaLog/config"
	"go.permanent.de/anaLog/logpoint"
	"go.permanent.de/anaLog/persistence/influx"
	"go.permanent.de/anaLog/persistence/mongo"
	"go.permanent.de/anaLog/state"
)

func Close() {
	if persistenceAdapter != nil {
		getAdapter().Close()
	}
}

var persistenceAdapter Adapter

var (
	NXlogpoint error = errors.New("LogPoint not found")
)

type Adapter interface {
	StorePoint(logpoint.LogPoint) error
	StorePoints(...logpoint.LogPoint) error
	GetRecurring() (map[string]map[string]map[string]logpoint.LogPoint, error)
	StoreAnalysisRCserial([]byte) error
	GetLatestAnalysisRCserial() ([]byte, error)
	GetEndByStart(logpoint.LogPoint) (logpoint.LogPoint, error)
	GetLastBegin(string) (logpoint.LogPoint, error)
	Find(task, host, state, rawRegex string, timeRangeGTE, timeRangeLTE time.Time, n uint) ([]logpoint.LogPoint, error)
	Close() error
}

func getAdapter() Adapter {
	if persistenceAdapter == nil {
		if config.AnaLog.UseInflux {
			persistenceAdapter = influx.GetAdapter()
		} else if config.AnaLog.UseMongo {
			persistenceAdapter = mongo.GetAdapter()
		}
	}
	return persistenceAdapter
}

func StorePoint(lp logpoint.LogPoint) error {
	return getAdapter().StorePoint(lp)
}

func StorePoints(lps ...logpoint.LogPoint) error {
	return getAdapter().StorePoints(lps...)
}

func GetRecurring() (map[string]map[string]map[string]logpoint.LogPoint, error) {
	//--------------------task----identifier---state------------------
	return getAdapter().GetRecurring()
}

func StoreAnalysisRCserial(lrc []byte) error {
	return getAdapter().StoreAnalysisRCserial(lrc)
}

func GetLatestAnalysisRCserial() ([]byte, error) {
	return getAdapter().GetLatestAnalysisRCserial()
}

func GetEndByStart(begin logpoint.LogPoint) (logpoint.LogPoint, error) {
	ret, err := getAdapter().GetEndByStart(begin)
	if err != nil {
		return ret, nil
	} else if ret.State == state.Unknown && ret.RunId != begin.RunId {
		return ret, NXlogpoint
	}
	return ret, err
}

func GetLastBegin(taskname string) (logpoint.LogPoint, error) {
	return getAdapter().GetLastBegin(taskname)
}

func Find(task, host, state, rawRegex string, timeRangeGTE, timeRangeLTE time.Time, n uint) ([]logpoint.LogPoint, error) {
	lps, err := getAdapter().Find(task, host, state, rawRegex, timeRangeGTE, timeRangeLTE, n)
	if err != nil {
		return nil, err
	}
	sort.Sort(logpoint.ByTime(lps))
	return lps, err
}
