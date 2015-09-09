package anaLog

import (
	"strconv"
	"time"

	"go.iondynamics.net/iDhelper/randGen"
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"

	"go.permanent.de/anaLog/anaLog/heartbeat"
	"go.permanent.de/anaLog/anaLog/logpoint"
	"go.permanent.de/anaLog/anaLog/mode"
	"go.permanent.de/anaLog/anaLog/persistence"
	"go.permanent.de/anaLog/anaLog/scheduler"
	"go.permanent.de/anaLog/anaLog/state"
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
	heartbeat.Create(heartbeat.Heartbeat{lp, "heartbeat"})
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
		idl.Err("Recurring task unsuccessful: "+lp.Task, lp)
	}
	heartbeat.Exit(heartbeat.Heartbeat{lp, "exit"})
	return persistence.StorePoint(lp)
}

func PushRecurringHeartbeat(host, task, identifier, subtask string) error {
	hb := heartbeat.Heartbeat{
		LogPoint: logpoint.LogPoint{
			RunId:    newRunId(),
			Task:     task,
			Host:     host,
			Mode:     mode.Recurring,
			Priority: priority.Informational,
			State:    state.Started,
			Time:     time.Now(),
		},
		Subtask: subtask,
	}
	heartbeat.Ping(hb)
	return nil
}
