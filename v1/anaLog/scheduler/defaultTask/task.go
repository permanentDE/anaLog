package defaultTask

import (
	"time"

	idl "go.iondynamics.net/iDlogger"

	"go.permanent.de/anaLog/v1/anaLog/analysis"
	"go.permanent.de/anaLog/v1/anaLog/scheduler"
)

func init() {
	ch := make(chan time.Time)
	scheduler.Register(ch)
	go loop(ch)
}

func loop(ch chan time.Time) {
	for {
		<-ch
		idl.Info("default task: tick")
		err := analysis.CheckRecurringFluctuation()
		if err != nil {
			idl.Crit("scheduled analysis.Check failed: ", err)
		}
	}
}
