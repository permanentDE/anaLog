package logpoint

import (
	"time"

	"go.iondynamics.net/iDlogger/priority"

	"go.permanent.de/anaLog/anaLog/mode"
	"go.permanent.de/anaLog/anaLog/state"
)

type LogPoint struct {
	RunId    string
	Task     string
	Host     string
	Mode     mode.Mode
	Priority priority.Priority
	State    state.State
	Time     time.Time
	Message  string
	Raw      string
	Data     map[string]interface{}
}
