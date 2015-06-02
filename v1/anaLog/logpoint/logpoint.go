package logpoint

import (
	"time"

	"go.iondynamics.net/iDlogger/priority"

	"go.permanent.de/anaLog/v1/anaLog/mode"
	"go.permanent.de/anaLog/v1/anaLog/state"
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
