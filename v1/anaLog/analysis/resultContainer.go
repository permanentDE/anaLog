package analysis

import (
	"encoding/json"
	"sync"
	"time"

	"go.permanent.de/anaLog/v1/anaLog/persistence"
)

type Result struct {
	Avg    time.Duration
	StdDev time.Duration

	AvgQr    time.Duration
	StdDevQr time.Duration

	IntervalAvg    time.Duration
	IntervalStdDev time.Duration
}

type ResultContainer struct {
	mu sync.RWMutex
	m  map[string]Result
}

func NewResultContainer() *ResultContainer {
	rc := new(ResultContainer)
	rc.m = make(map[string]Result)
	return rc
}

func (rc *ResultContainer) Set(key string, r Result) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.m[key] = r
}

func (rc *ResultContainer) Get(key string) Result {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return rc.m[key]
}

func (rc *ResultContainer) Range(f func(string, Result)) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	for key, val := range rc.m {
		f(key, val)
	}
}

var latestRcCache *ResultContainer = NewResultContainer()

func (rc *ResultContainer) Store() error {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if latestRcCache != nil {
		latestRcCache.mu.Lock()
		defer latestRcCache.mu.Unlock()
	}

	latestRcCache = rc
	serial, err := json.Marshal(rc.m)
	if err != nil {
		return err
	}
	return persistence.StoreAnalysisRCserial(serial)
}

func (rc *ResultContainer) LoadLatest() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if latestRcCache != nil && len(latestRcCache.m) > 0 {
		latestRcCache.mu.RLock()
		defer latestRcCache.mu.RUnlock()

		rc.m = latestRcCache.m
		return nil
	}

	lrcSerial, err := persistence.GetLatestAnalysisRCserial()
	if err == nil {
		err = json.Unmarshal(lrcSerial, &latestRcCache.m)
		if err != nil {
			return err
		}
		rc.m = latestRcCache.m
	}
	return err
}
