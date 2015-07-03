package iDlogger

import (
	"reflect"
)

type LogError struct {
	Event *Event
	Err   error
}

func (err LogError) Error() string {
	return err.Err.Error()
}

func IsLogError(err interface{}) bool {
	if reflect.TypeOf(err).String() == "iDlogger.LogError" {
		return true
	} else {
		return false
	}
}
