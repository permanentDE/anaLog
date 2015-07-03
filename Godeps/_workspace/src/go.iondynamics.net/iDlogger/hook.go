package iDlogger

import (
	"io"

	"go.iondynamics.net/iDlogger/priority"
)

type Hook interface {
	Priorities() []priority.Priority
	Fire(*Event) error
}

type StdHook struct {
	p []priority.Priority
	w io.Writer
	f Formatter
}

func (sh *StdHook) Fire(e *Event) error {
	byt, err := sh.f.Format(e)
	if err == nil {
		_, err = sh.w.Write(*byt)
	}
	return err
}

func (sh *StdHook) Priorities() []priority.Priority {
	return sh.p
}

func (sh *StdHook) SetPriorities(p []priority.Priority) {
	sh.p = p
}

func (sh *StdHook) SetWriter(w io.Writer) {
	sh.w = w
}

func (sh *StdHook) SetFormatter(f Formatter) {
	sh.f = f
}
