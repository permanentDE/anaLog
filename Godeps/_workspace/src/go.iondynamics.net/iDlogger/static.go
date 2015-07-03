package iDlogger

import (
	"fmt"
	"os"
	"time"

	"go.iondynamics.net/iDlogger/priority"
)

var (
	std = New()
)

func StandardLogger() *Logger {
	return std
}

func Wait() {
	std.wg.Wait()
}

func SetPrefix(prefix string) {
	std.prefix = prefix
}

func SetErrCallback(errCallback func(error)) {
	std.errCallback = errCallback
}

func AddHook(h Hook) {
	std.mu.Lock()
	defer std.mu.Unlock()

	for _, prio := range h.Priorities() {
		std.priorityHooks[prio] = append(std.priorityHooks[prio], h)
	}
}

func Log(e *Event) {
	std.wg.Add(1)
	if std.Async {
		go std.dispatch(e)
	} else {
		std.dispatch(e)
		std.wg.Wait()
	}
}

func Debug(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Debugging, fmt.Sprint(entry...)})
}

func Info(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Informational, fmt.Sprint(entry...)})
}

func Notice(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Notice, fmt.Sprint(entry...)})
}

func Warn(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Warning, fmt.Sprint(entry...)})
}

func Err(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Error, fmt.Sprint(entry...)})
}

func Crit(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Critical, fmt.Sprint(entry...)})
}

func Alert(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Alert, fmt.Sprint(entry...)})
}

func Emerg(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Emergency, fmt.Sprint(entry...)})
	std.Wait()
	os.Exit(1)
}

func Panic(entry ...interface{}) {
	std.Log(&Event{std, map[string]interface{}{}, time.Now(), priority.Emergency, fmt.Sprint(entry...)})
	std.Wait()
	panic(fmt.Sprint(entry...))
}
