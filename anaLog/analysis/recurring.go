package analysis

import (
	"errors"
	"sync"
	"time"

	"go.permanent.de/anaLog/anaLog/logpoint"
	"go.permanent.de/anaLog/anaLog/persistence"
	"go.permanent.de/anaLog/anaLog/state"
)

var (
	NoRecurringData error = errors.New("No recurring data (yet)")
)

func GetRecurringResultContainer() (*ResultContainer, error) {
	taskAnalysis := NewResultContainer()
	table, err := persistence.GetRecurring()
	if err != nil {
		return taskAnalysis, err
	}

	taskDurations := make(map[string][]time.Duration)
	taskBeginning := make(map[string][]time.Time)

	for task, idstatemap := range table {
	RunIterator:
		for id, statemap := range idstatemap {
			_ = id
			var startLp logpoint.LogPoint
			var endLp logpoint.LogPoint

			var foundStart bool
			var foundEnd bool

		StartEndLookup:
			for _, lp := range statemap {
				if !foundStart && lp.State == state.Started {
					previous := taskBeginning[task]
					previous = append(previous, lp.Time)
					taskBeginning[task] = previous

					startLp = lp
					foundStart = true
				}
				if !foundEnd && lp.State != state.Started && lp.State != state.Running && lp.State != state.Unknown {
					endLp = lp
					foundEnd = true
				}

				if foundStart && foundEnd {
					break StartEndLookup //short circuit
				}
			}

			if !(foundStart && foundEnd) {
				continue RunIterator //incomplete tuple
			}

			duration := endLp.Time.Sub(startLp.Time)
			previous := taskDurations[task]
			previous = append(previous, duration)
			taskDurations[task] = previous
		}
	}

	outerWg := sync.WaitGroup{}

	for task, durations := range taskDurations {
		outerWg.Add(1)

		go func(name string, durList []time.Duration) {
			innerWg := sync.WaitGroup{}
			result := &Result{}

			innerWg.Add(1)
			go func(ds []time.Duration) {
				result.Avg = Avg(ds)
				innerWg.Done()
			}(durList)

			innerWg.Add(1)
			go func(ds []time.Duration) {
				result.StdDev = StdDev(ds, -1)
				innerWg.Done()
			}(durList)

			innerWg.Add(3)
			go func(times []time.Time) {
				durs := DurationsBetween(times)

				go func() {
					result.IntervalAvg = Avg(durs)
					innerWg.Done()
				}()
				go func() {
					result.IntervalStdDev = StdDev(durs, -1)
					innerWg.Done()
				}()
				innerWg.Done()
			}(taskBeginning[name])

			durs := append([]time.Duration{}, durList...)
			qrDur := QuartileReduce(durs)

			innerWg.Add(1)
			go func(ds []time.Duration) {
				result.AvgQr = Avg(ds)
				innerWg.Done()
			}(qrDur)

			innerWg.Add(1)
			go func(ds []time.Duration) {
				result.StdDevQr = StdDev(ds, -1)
				innerWg.Done()
			}(qrDur)

			innerWg.Wait()
			taskAnalysis.Set(name, *result)
			outerWg.Done()
		}(task, durations)

	}
	outerWg.Wait()
	return taskAnalysis, nil
}

func RecurringExpectedAfter(lp logpoint.LogPoint) (time.Duration, error) {
	lastRc := NewResultContainer()
	err := lastRc.LoadLatest()
	if err != nil {
		if lastRc == nil {
			return 0, err
		}
	}

	res := lastRc.Get(lp.Task)
	if res.Avg < 1 {
		return 0, NoRecurringData
	}

	return res.AvgQr + res.StdDev + res.StdDevQr, nil
}
