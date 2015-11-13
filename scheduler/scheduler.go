package scheduler

import (
	"time"

	idl "go.iondynamics.net/iDlogger"

	"go.permanent.de/anaLog/alertBlocker"
	"go.permanent.de/anaLog/analysis"
	"go.permanent.de/anaLog/config"
	"go.permanent.de/anaLog/heartbeat"
	"go.permanent.de/anaLog/logpoint"
)

var registeredChannels []chan<- time.Time
var gracePeriod time.Duration

func Start() {
	if config.AnaLog.SchedulerInterval == "" {
		idl.Notice("Scheduling disabled")
		return
	}
	dur, err := time.ParseDuration(config.AnaLog.SchedulerInterval)
	if err != nil {
		idl.Emerg("Invalid configuration: AnaLog.SchedulerInterval")
	}

	if config.AnaLog.GracePeriod != "" {
		gracePeriod, err = time.ParseDuration(config.AnaLog.GracePeriod)
		if err != nil {
			idl.Emerg("Invalid configuration: AnaLog.GracePeriod")
		}
	}

	go func() {
		pingAll(time.Now())
	}()

	c := time.Tick(dur)
	go loop(c)
}

func StartIn(dur time.Duration) {
	go func(d time.Duration) {
		<-time.After(dur)
		Start()
	}(dur)
}

func loop(ch <-chan time.Time) {
	for now := range ch {
		pingAll(now)
	}
}

func pingAll(t time.Time) {
	for _, channel := range registeredChannels {
		channel <- t
	}
}

func Register(channel chan<- time.Time) {
	registeredChannels = append(registeredChannels, channel)
}

func RecurringTaskIncoming(begin logpoint.LogPoint) {
	dur, err := analysis.RecurringExpectedAfter(begin)
	if !alertBlocker.IsUnknown(begin.Task, "not recurred") {
		alertBlocker.Resolved(begin.Task, "not recurred")
	}

	if err == analysis.NoRecurringData {
		idl.Notice(`Skipping analysis scheduling of recurring task "` + begin.Task + `" due to missing data`)
		return
	}
	if err != nil {
		idl.Crit("Failed scheduling for analysis of recurring task ", err, begin)
	}

	<-time.After(dur + gracePeriod)
	heartbeat.Wait(heartbeat.Heartbeat{LogPoint: begin})
	err = analysis.CheckRecurredTaskEnd(begin)
	if err != nil {
		idl.Crit("Failed analysis of recurring task ", err, begin)
	}
}
