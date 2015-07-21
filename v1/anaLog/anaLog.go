package anaLog

import (
	"strconv"
	"time"

	"go.iondynamics.net/iDhelper/randGen"
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"

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
