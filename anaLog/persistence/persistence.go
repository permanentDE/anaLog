package persistence

import (
	"errors"

	"go.permanent.de/anaLog/anaLog/logpoint"
	"go.permanent.de/anaLog/anaLog/persistence/influx"
	"go.permanent.de/anaLog/anaLog/persistence/mongo"
	"go.permanent.de/anaLog/anaLog/state"
	"go.permanent.de/anaLog/config"
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
	/*	GetPoint(string) (logpoint.LogPoint, error)
		UpsertPoint(logpoint.LogPoint) error*/
	StorePoint(logpoint.LogPoint) error
	StorePoints(...logpoint.LogPoint) error
	GetRecurring() (map[string]map[string]map[string]logpoint.LogPoint, error)
	StoreAnalysisRCserial([]byte) error
	GetLatestAnalysisRCserial() ([]byte, error)
	GetEndByStart(logpoint.LogPoint) (logpoint.LogPoint, error)
	GetLastBegin(string) (logpoint.LogPoint, error)
	Close()
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

/*func GetPoint(identifier string) (logpoint.LogPoint, error) {
	return getAdapter().GetPoint(identifier)
}

func UpsertPoint(lp logpoint.LogPoint) error {
	return getAdapter().UpsertPoint(lp)
}
*/
