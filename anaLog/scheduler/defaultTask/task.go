package defaultTask

import (
	"time"

	idl "go.iondynamics.net/iDlogger"

	"go.permanent.de/anaLog/anaLog/analysis"
	"go.permanent.de/anaLog/anaLog/scheduler"
	"go.permanent.de/anaLog/config"
)

func init() {
	flucCh := make(chan time.Time)
	scheduler.Register(flucCh)
	go fluctuationLoop(flucCh)

	recurCh := make(chan time.Time)
	scheduler.Register(recurCh)
	go recurringBeginWatcher(recurCh)

}

func fluctuationLoop(ch chan time.Time) {
	for {
		<-ch
		err := analysis.CheckRecurringFluctuation()
		if err != nil {
			idl.Crit("scheduled analysis.Check failed: ", err)
		}
	}
}

func recurringBeginWatcher(ch chan time.Time) {
	schedulerInterval, _ := time.ParseDuration(config.AnaLog.SchedulerInterval)

	for {
		<-ch
		rc := analysis.NewResultContainer()
		err := rc.LoadLatest()
		if err != nil {
			idl.Crit("Failed scheduling of task begin analysis: ", err)
		}

		rc.Range(createAnalyzer(time.Now().Add(schedulerInterval)))

	}
}

func createAnalyzer(validUntil time.Time) func(string, analysis.Result) {
	vu := validUntil
	fn := func(t string, r analysis.Result) {
		go func(task string, res analysis.Result, valid time.Time) {
			for {
				if time.Since(valid) > 0 || res.IntervalAvg == 0 {
					return
				}
				analysis.CheckRecurredTaskBegin(task)
				<-time.After(res.IntervalAvg + res.IntervalStdDev)
			}
		}(t, r, vu)
	}

	return fn
}
