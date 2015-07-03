package priority

import (
	"strings"
)

type Priority uint8

const (
	Invalid Priority = Priority(^uint8(0)) //Maximum of uint8
)

const (
	Emergency Priority = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Informational
	Debugging
)

func (priority Priority) String() string {
	switch priority {
	case Emergency:
		return "emerg"
	case Alert:
		return "alert"
	case Critical:
		return "crit"
	case Error:
		return "err"
	case Warning:
		return "warn"
	case Notice:
		return "notice"
	case Informational:
		return "info"
	case Debugging:
		return "debug"
	}

	return "unknown"
}

func Atos(a string) Priority {
	switch strings.ToLower(strings.TrimSpace(a)) {

	case "emerg", "emergency":
		return Emergency

	case "alert":
		return Alert

	case "crit", "critical":
		return Critical

	case "err", "error":
		return Error

	case "warn", "warning":
		return Warning

	case "notice":
		return Notice

	case "info", "informational":
		return Informational

	case "debug", "debugging":
		return Debugging
	}

	return Invalid

}

var allPriorities = []Priority{
	Debugging,
	Informational,
	Notice,
	Warning,
	Error,
	Critical,
	Alert,
	Emergency,
}

func Threshold(p Priority) []Priority {
	for i := range allPriorities {
		if allPriorities[i] == p {
			return allPriorities[i:]
		}
	}
	return []Priority{}
}
