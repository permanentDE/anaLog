package logpoint

import (
	"time"

	"go.iondynamics.net/iDlogger/priority"

	"go.permanent.de/anaLog/mode"
	"go.permanent.de/anaLog/state"
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

type ByTime []LogPoint

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }
