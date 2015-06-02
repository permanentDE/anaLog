package anaLog

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"go.iondynamics.net/iDhelper/randGen"
	"go.iondynamics.net/iDlogger/priority"

	"go.permanent.de/anaLog/v1/anaLog/analysis"
	"go.permanent.de/anaLog/v1/anaLog/logpoint"
	"go.permanent.de/anaLog/v1/anaLog/mode"
	"go.permanent.de/anaLog/v1/anaLog/persistence"
	"go.permanent.de/anaLog/v1/anaLog/state"
)

func newRunId() string {
	return strconv.Itoa(int(time.Now().UnixNano())) + "_" + randGen.String(64)
}

func GenerateSampleData() interface{} {
	var points []logpoint.LogPoint

	rand.Seed(time.Now().UnixNano())
	//return randGen.String(123)
	day := 24 * time.Hour
	monthStart := time.Now().Add(-1 * day * 1000)
	for i := 0; i < 1000; i++ {
		taskStart := monthStart.Add(day * time.Duration(i))
		taskEnd := taskStart.Add(500 * time.Second).Add(time.Duration(rand.Intn(100)) * time.Second)
		startLp := logpoint.LogPoint{
			RunId:    newRunId(),
			Task:     "sampleData",
			Host:     "example.permanent.de",
			Mode:     mode.Recurring,
			Priority: priority.Informational,
			State:    state.Started,
			Time:     taskStart,
		}
		endLp := startLp
		endLp.Time = taskEnd
		endLp.State = state.OK
		points = append(points, startLp, endLp)
	}

	return persistence.StorePoints(points...)
}

func PushRecurringBegin(task, host string) (string, error) {
	lp := logpoint.LogPoint{
		RunId:    newRunId(),
		Task:     task,
		Host:     host,
		Mode:     mode.Recurring,
		Priority: priority.Informational,
		State:    state.Started,
		Time:     time.Now(),
	}
	return lp.RunId, persistence.StorePoint(lp)
}

func PushRecurringEnd(task, host, identifier, stateStr, requestBody string) error {
	lp := logpoint.LogPoint{
		RunId:    identifier,
		Task:     task,
		Host:     host,
		Mode:     mode.Recurring,
		Priority: priority.Informational,
		State:    state.Atos(stateStr),
		Time:     time.Now(),
		Raw:      requestBody,
	}
	return persistence.StorePoint(lp)
}

func AnalyzeRecurring() interface{} {
	//return GenerateSampleData()//don't execute this more than once

	table, err := persistence.GetRecurring()
	if err != nil {
		return err
	}

	taskDurations := make(map[string][]time.Duration)
	taskAnalysis := make(map[string]map[string]time.Duration)

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
		taskAnalysis[task] = make(map[string]time.Duration)
	}

	outerWg := sync.WaitGroup{}
	for task, durations := range taskDurations {
		outerWg.Add(1)

		go func() {
			precision := time.Millisecond

			internalWg := sync.WaitGroup{}
			internalWg.Add(4)
			go func() {
				taskAnalysis[task]["avg"] = analysis.Avg(durations)
				internalWg.Done()
			}()
			go func() {
				taskAnalysis[task]["stdDeviation"] = analysis.StdDev(durations, precision)
				internalWg.Done()
			}()

			qrDur := analysis.QuartileReduce(durations)

			go func() {
				taskAnalysis[task]["avgQr"] = analysis.Avg(qrDur)
				internalWg.Done()
			}()
			go func() {
				taskAnalysis[task]["stdDeviationQr"] = analysis.StdDev(qrDur, precision)
				internalWg.Done()
			}()
			outerWg.Done()
		}()

	}
	outerWg.Wait()

	return taskAnalysis
}
