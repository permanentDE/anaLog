package analysis

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.permanent.de/anaLog/persistence"
)

type Result struct {
	Avg    time.Duration
	StdDev time.Duration

	AvgQr    time.Duration
	StdDevQr time.Duration

	IntervalAvg    time.Duration
	IntervalStdDev time.Duration
}

//String returns a textual representation of the result
func (r Result) String() string {
	return fmt.Sprint("Avg: ", r.Avg, ", StdDev: ", r.StdDev, ", AvgQr: ", r.AvgQr, ", StdDevQr: ", r.StdDevQr, ", IntervalAvg: ", r.IntervalAvg, ", IntervalStdDev: ", r.IntervalStdDev)
}

type ResultContainer struct {
	mu sync.RWMutex
	m  map[string]Result
}

//NewResultContainer returns a initialized ResultContainer (Pointer)
func NewResultContainer() *ResultContainer {
	rc := new(ResultContainer)
	rc.m = make(map[string]Result)
	return rc
}

//Set saves the given Result with the given key
func (rc *ResultContainer) Set(key string, r Result) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.m[key] = r
}

//Get returns the Result associated with the given key
func (rc *ResultContainer) Get(key string) Result {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return rc.m[key]
}

//Range calls the given function for every task in the container
func (rc *ResultContainer) Range(f func(string, Result)) {
	rc.mu.RLock()
	copy := make(map[string]Result)
	for key, val := range rc.m {
		copy[key] = val
	}
	rc.mu.RUnlock()

	for key, val := range copy {
		f(key, val)
	}
}

var latestRcCache *ResultContainer = NewResultContainer()

//Store saves the ResultContainer in memory and stores a serialized version in persistence
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

//LoadLatest gets the latest data from memory or the persistence layer
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
