package iDlogger

import (
	"go.iondynamics.net/iDlogger/priority"
	"time"
)

type Event struct {
	Logger   *Logger
	Data     map[string]interface{}
	Time     time.Time
	Priority priority.Priority
	Message  string
}
