package state

import (
	"strings"
)

type State string

const (
	Started            State = "Started"
	Running            State = "Running"
	OK                 State = "OK"
	Failed             State = "Failed"
	CompletedWithError State = "CompletedWithError"
	Unknown            State = "Unknown"
)

func Atos(a string) State {

	switch strings.ToLower(a) {
	case "started":
		return Started
	case "running":
		return Running
	case "ok":
		return OK
	case "failed":
		return Failed
	case "completedwitherror":
		return CompletedWithError
	default:
		return Unknown
	}

}
