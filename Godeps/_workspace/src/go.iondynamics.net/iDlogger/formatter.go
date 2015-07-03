package iDlogger

import (
	"fmt"
)

type Formatter interface {
	Format(e *Event) (*[]byte, error)
}

type StdFormatter struct {
}

func (sf *StdFormatter) Format(e *Event) (*[]byte, error) {
	var buf []byte
	if len(e.Logger.prefix) > 0 {
		buf = append(buf, []byte("["+e.Logger.prefix+"]")...)
	}
	buf = append(buf, []byte("["+e.Priority.String()+"]["+e.Time.Format("2006-01-02 15:04:05")+"]\t"+e.Message)...)
	i := 1
	for k, v := range e.Data {
		if i == 1 && e.Message == "" {
			buf = append(buf, []byte(k+":"+fmt.Sprint(v))...)
		} else {
			buf = append(buf, []byte("\t"+k+":"+fmt.Sprint(v))...)
		}
		i++
	}
	buf = append(buf, []byte("\n")...)
	return &buf, nil
}
