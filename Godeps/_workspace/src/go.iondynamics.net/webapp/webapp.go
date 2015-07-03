package webapp

import (
	idl "go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"
	"net/http"
	"time"
)

var Std = idl.StandardLogger()

type Error struct {
	Error   error
	Message string
	Code    int
	Write   bool
}

func New(err error, message string, code int) *Error {
	return &Error{Error: err, Message: message, Code: code, Write: false}
}

func Write(err error, message string, code int) *Error {
	return &Error{Error: err, Message: message, Code: code, Write: true}
}

type Handler func(http.ResponseWriter, *http.Request) *Error

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		Std.Log(&idl.Event{Std, map[string]interface{}{"Error": e.Error, "Message": e.Message, "Code": e.Code}, time.Now(), priority.Error, e.Error.Error()})
		if e.Code != 0 && e.Write {
			http.Error(w, e.Message, e.Code)
		}
	}
}
