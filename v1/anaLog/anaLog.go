package anaLog

import (
	"math/rand"
	"strconv"
	"time"

	"go.iondynamics.net/iDhelper/randGen"
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"

	"go.permanent.de/anaLog/v1/anaLog/analysis"
	"go.permanent.de/anaLog/v1/anaLog/logpoint"
	"go.permanent.de/anaLog/v1/anaLog/mode"
	"go.permanent.de/anaLog/v1/anaLog/persistence"
	"go.permanent.de/anaLog/v1/anaLog/scheduler"
	"go.permanent.de/anaLog/v1/anaLog/state"
)

func newRunId() string {
	return strconv.Itoa(int(time.Now().UnixNano())) + "_" + randGen.String(64)
}

func Close() {
	persistence.Close()
}

func GenerateSampleData() interface{} {
	var points []logpoint.LogPoint

	rand.Seed(time.Now().UnixNano())
	day := 24 * time.Hour
	monthStart := time.Now().Add(-1 * day * 100)
	for i := 0; i < 100; i++ {
		taskStart := monthStart.Add(day * time.Duration(i))
		taskEnd := taskStart.Add(30 * time.Second).Add(time.Duration(rand.Intn(10)) * time.Second)
		startLp := logpoint.LogPoint{
			RunId:    newRunId(),
			Task:     "testData",
			Host:     "test.permanent.de",
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
	go scheduler.RecurringTaskIncoming(lp)
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
	if lp.State != state.OK {
		idl.Warn("Recurring task unsuccessful: "+lp.Task, lp)
	}
	return persistence.StorePoint(lp)
}

func AnalyzeRecurring() (*analysis.ResultContainer, error) {
	//return GenerateSampleData() //don't execute this more than once
	/*lp := logpoint.LogPoint{
		RunId:    "1435242237863904600_TY1IFnZMJgJ2oQ1cbCZ1Noc7srSpTk2GqvWvyCFRkjiH9KtBLtLk21TResavgeAr",
		Task:     "testData",
		Host:     "test.permanent.de",
		Mode:     mode.Recurring,
		Priority: priority.Informational,
		State:    state.Started,
		Time:     time.Now(),
	}

	go scheduler.RecurringTaskIncoming(lp)*/

	analysis.CheckRecurredTaskBegin("testData")

	return analysis.GetRecurringResultContainer()
}

func diff(X, Y []time.Duration) []time.Duration {
	m := make(map[time.Duration]int)

	for _, y := range Y {
		m[y]++
	}
	var ret []time.Duration
	for _, x := range X {
		if m[x] > 0 {
			m[x]--
			continue
		}
		ret = append(ret, x)
	}

	return ret
}
