package scheduler

import (
	"time"

	idl "go.iondynamics.net/iDlogger"

	"go.permanent.de/anaLog/anaLog/analysis"
	"go.permanent.de/anaLog/anaLog/heartbeat"
	"go.permanent.de/anaLog/anaLog/logpoint"
	"go.permanent.de/anaLog/config"
)

var registeredChannels []chan<- time.Time
var GracePeriod time.Duration

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
		GracePeriod, err = time.ParseDuration(config.AnaLog.GracePeriod)
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
	if err == analysis.NoRecurringData {
		idl.Notice(`Skipping analysis scheduling of recurring task "` + begin.Task + `" due to missing data`)
		return
	}
	if err != nil {
		idl.Crit("Failed scheduling for analysis of recurring task ", err, begin)
	}

	<-time.After(dur + GracePeriod)
	<-heartbeat.StillAlive(begin.Task, begin.RunId)
	err = analysis.CheckRecurredTaskEnd(begin)
	if err != nil {
		idl.Crit("Failed analysis of recurring task ", err, begin)
	}
}
