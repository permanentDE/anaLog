package iDlogger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.iondynamics.net/iDlogger/priority"
)

type Logger struct {
	Async         bool
	mu            sync.RWMutex
	wg            sync.WaitGroup
	prefix        string
	flag          uint8
	errCallback   func(error)
	priorityHooks map[priority.Priority][]Hook
}

func New() *Logger {
	var sf *StdFormatter
	stdOut := new(StdHook)
	stdErr := new(StdHook)

	stdOut.SetFormatter(sf)
	stdErr.SetFormatter(sf)

	stdOut.SetWriter(os.Stdout)
	stdErr.SetWriter(os.Stderr)

	stdOut.SetPriorities([]priority.Priority{
		priority.Debugging,
		priority.Informational,
		priority.Notice,
	})

	stdErr.SetPriorities(priority.Threshold(priority.Warning))

	log := &Logger{
		false,
		sync.RWMutex{},
		sync.WaitGroup{},
		"",
		0,
		nil,
		map[priority.Priority][]Hook{
			priority.Emergency:     []Hook{},
			priority.Alert:         []Hook{},
			priority.Critical:      []Hook{},
			priority.Error:         []Hook{},
			priority.Warning:       []Hook{},
			priority.Notice:        []Hook{},
			priority.Informational: []Hook{},
			priority.Debugging:     []Hook{},
		},
	}

	log.AddHook(stdOut)
	log.AddHook(stdErr)

	return log
}

func (l *Logger) dispatch(e *Event) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, h := range l.priorityHooks[e.Priority] {
		err := h.Fire(e)
		if err != nil && l.errCallback != nil {
			l.errCallback(err)
		}
	}
	l.wg.Done()
}

func (l *Logger) Wait() {
	l.wg.Wait()
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) SetErrCallback(errCallback func(error)) {
	l.errCallback = errCallback
}

func (l *Logger) AddHook(h Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, prio := range h.Priorities() {
		l.priorityHooks[prio] = append(l.priorityHooks[prio], h)
	}
}

func (l *Logger) Log(e *Event) {
	l.wg.Add(1)
	if l.Async {
		go l.dispatch(e)
	} else {
		l.dispatch(e)
		l.wg.Wait()
	}
}

func (l *Logger) Debug(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Debugging, fmt.Sprint(entry...)})
}

func (l *Logger) Info(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Informational, fmt.Sprint(entry...)})
}

func (l *Logger) Notice(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Notice, fmt.Sprint(entry...)})
}

func (l *Logger) Warn(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Warning, fmt.Sprint(entry...)})
}

func (l *Logger) Err(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Error, fmt.Sprint(entry...)})
}

func (l *Logger) Crit(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Critical, fmt.Sprint(entry...)})
}

func (l *Logger) Alert(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Alert, fmt.Sprint(entry...)})
}

func (l *Logger) Emerg(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Emergency, fmt.Sprint(entry...)})
	l.Wait()
	os.Exit(1)
}

func (l *Logger) Panic(entry ...interface{}) {
	l.Log(&Event{l, map[string]interface{}{}, time.Now(), priority.Emergency, fmt.Sprint(entry...)})
	l.Wait()
	panic(fmt.Sprint(entry...))
}
