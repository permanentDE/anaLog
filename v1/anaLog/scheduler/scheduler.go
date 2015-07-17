package scheduler

import (
	"time"

	idl "go.iondynamics.net/iDlogger"

	"go.permanent.de/anaLog/v1/anaLog/analysis"
	"go.permanent.de/anaLog/v1/anaLog/logpoint"
	"go.permanent.de/anaLog/v1/config"
)

var registeredChannels []chan<- time.Time

func Start() {
	if config.AnaLog.SchedulerInterval == "" {
		idl.Notice("Scheduling disabled")
		return
	}
	dur, err := time.ParseDuration(config.AnaLog.SchedulerInterval)
	if err != nil {
		idl.Emerg("Invalid configuration: AnaLog.SchedulerInterval")
	}

	go func() {
		pingAll(time.Now())
	}()

	c := time.Tick(dur)
	go loop(c)
	go RecurringBeginWatcher()
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

	<-time.After(dur)
	err = analysis.CheckRecurredTaskEnd(begin)
	if err != nil {
		idl.Crit("Failed analysis of recurring task ", err, begin)
	}
}

func RecurringBeginWatcher() {
	rc := analysis.NewResultContainer()
	err := rc.LoadLatest()
	if err != nil {
		idl.Crit("Failed scheduling of task begin analysis ", err)
	}

	analyzer := func(t string, r analysis.Result) {
		go func(task string, res analysis.Result) {
			for {
				analysis.CheckRecurredTaskBegin(task)
				<-time.After(res.IntervalAvg + res.IntervalStdDev)
			}
		}(t, r)
	}

	rc.Range(analyzer)
}
