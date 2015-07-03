package iDnegroniLog

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"go.iondynamics.net/iDlogger"
	"go.iondynamics.net/iDlogger/priority"
)

// Middleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type Middleware struct {
	Logger *iDlogger.Logger

	Priority priority.Priority

	LogPanicsWithPriority priority.Priority
	Stack2Http            bool
}

// NewMiddleware returns a new *Middleware, yay!
func NewMiddleware(log *iDlogger.Logger) *Middleware {
	return NewCustomMiddleware(log, priority.Informational, priority.Emergency, true)
}

func NewCustomMiddleware(log *iDlogger.Logger, prio, logPanicsWithPriority priority.Priority, stack2http bool) *Middleware {
	return &Middleware{Logger: log, Priority: prio, LogPanicsWithPriority: logPanicsWithPriority, Stack2Http: stack2http}
}

func (l *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 4*1024)
			stack = stack[:runtime.Stack(stack, false)]
			l.Logger.Log(&iDlogger.Event{
				l.Logger,
				map[string]interface{}{
					"panic":       err,
					"status":      http.StatusInternalServerError,
					"method":      r.Method,
					"request":     r.RequestURI,
					"remote":      r.RemoteAddr,
					"text_status": http.StatusText(http.StatusInternalServerError),
					"stack":       string(stack),
				},
				time.Now(),
				l.LogPanicsWithPriority,
				"PANIC while handling request",
			})

			http.Error(rw, "internal server error", http.StatusInternalServerError)

			if l.Stack2Http {
				fmt.Fprintf(rw, "PANIC: %s\n%s", err, stack)
			}

			l.Logger.Wait()
		}
	}()

	start := time.Now()
	startEvent := &iDlogger.Event{
		l.Logger,
		map[string]interface{}{
			"method":  r.Method,
			"request": r.RequestURI,
			"remote":  r.RemoteAddr,
		},
		time.Now(),
		l.Priority,
		"started handling request",
	}
	l.Logger.Log(startEvent)

	next(rw, r)

	latency := time.Since(start)
	res := rw.(negroni.ResponseWriter)

	completedEvent := &iDlogger.Event{
		l.Logger,
		map[string]interface{}{
			"status":      strconv.Itoa(res.Status()),
			"method":      r.Method,
			"request":     r.RequestURI,
			"remote":      r.RemoteAddr,
			"text_status": http.StatusText(res.Status()),
			"took":        latency.String(),
		},
		time.Now(),
		l.Priority,
		"completed handling request",
	}
	l.Logger.Log(completedEvent)
}
