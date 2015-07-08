package analysis

import (
	"math"
	"time"

	idl "go.iondynamics.net/iDlogger"

	"go.permanent.de/anaLog/v1/anaLog/logpoint"
	"go.permanent.de/anaLog/v1/anaLog/nagios"
	"go.permanent.de/anaLog/v1/anaLog/persistence"
	"go.permanent.de/anaLog/v1/anaLog/state"
)

func CheckRecurringFluctuation() error {
	lastRc := NewResultContainer()
	err := lastRc.LoadLatest()
	if err != nil {
		if lastRc == nil {
			negativeCheck("internal error")
			return err
		}
	}

	currentRc, err := GetRecurringResultContainer()
	if err != nil {
		negativeCheck("internal error")
		return err
	}
	err = currentRc.Store()

	currentRc.mu.RLock()
	defer currentRc.mu.RUnlock()
	lastRc.mu.RLock()
	defer lastRc.mu.RUnlock()

	warnPre := `Analysis of task "`
	warnMid := `" found fluctuation in `

	for task, curRes := range currentRc.m {
		lastRes, ok := lastRc.m[task]
		if !ok {
			//new task - no past data
			continue
		}

		avgDiff := int64(math.Abs(float64(curRes.Avg) - float64(lastRes.Avg)))
		if avgDiff > int64(lastRes.StdDevQr) {
			negativeCheck(warnPre + task + warnMid + "Average")
			idl.Warn(warnPre+task+warnMid+"Average", time.Duration(avgDiff), lastRes, curRes)
		}

		stdDevDiff := int64(math.Abs(float64(curRes.StdDev) - float64(lastRes.StdDev)))
		if stdDevDiff > int64(lastRes.StdDevQr) {
			negativeCheck(warnPre + task + warnMid + "StandardDeviation")
			idl.Warn(warnPre+task+warnMid+"StandardDeviation", time.Duration(stdDevDiff), lastRes, curRes)
		}

		avgDiffQr := int64(math.Abs(float64(curRes.AvgQr) - float64(lastRes.AvgQr)))
		if avgDiffQr > int64(lastRes.StdDevQr) {
			negativeCheck(warnPre + task + warnMid + "Average (quartile reduced)")
			idl.Warn(warnPre+task+warnMid+"Average (quartile reduced)", time.Duration(avgDiffQr), lastRes, curRes)
		}

		stdDevQrDiff := int64(math.Abs(float64(curRes.StdDevQr) - float64(lastRes.StdDevQr)))
		if stdDevQrDiff > int64(lastRes.StdDevQr) {
			negativeCheck(warnPre + task + warnMid + "StandardDeviation (quartile reduced)")
			idl.Warn(warnPre+task+warnMid+"StandardDeviation (quartile reduced)", time.Duration(stdDevQrDiff), lastRes, curRes)
		}

		invAvgDiff := int64(math.Abs(float64(curRes.IntervalAvg) - float64(lastRes.IntervalAvg)))
		if invAvgDiff > int64(lastRes.IntervalStdDev) {
			negativeCheck(warnPre + task + warnMid + "IntervalAverage")
			idl.Warn(warnPre+task+warnMid+"IntervalAverage", time.Duration(invAvgDiff), lastRes, curRes)
		}

		invStdDevDiff := int64(math.Abs(float64(curRes.IntervalStdDev) - float64(lastRes.IntervalStdDev)))
		if invStdDevDiff > int64(lastRes.IntervalStdDev) {
			negativeCheck(warnPre + task + warnMid + "IntervalStandardDeviation")
			idl.Warn(warnPre+task+warnMid+"IntervalStandardDeviation", time.Duration(invStdDevDiff), lastRes, curRes)
		}
	}

	return err
}

func CheckRecurredTaskEnd(begin logpoint.LogPoint) error {
	end, err := persistence.GetEndByStart(begin)
	if err == persistence.NXlogpoint {
		idl.Warn(`Task "`+begin.Task+`" not finished`, begin)
		return nil
	}

	if err != nil {
		negativeCheck("internal error")
		return err
	}

	if end.State != state.OK {
		idl.Warn(`Task "`+begin.Task+`" unsuccessful`, end)
	}

	idl.Debug("task successfully checked", end)

	return nil

}

func CheckRecurredTaskBegin(name string) error {
	lastRc := NewResultContainer()
	err := lastRc.LoadLatest()
	if err != nil {
		if lastRc == nil {
			negativeCheck("internal error")
			return err
		}
	}

	lastBegin, err := persistence.GetLastBegin(name)
	if err != nil {
		negativeCheck("internal error")
		return err
	}

	res := lastRc.Get(name)
	if res.IntervalAvg < 1 {
		return nil
		//no interval data
	}

	if time.Now().Sub(lastBegin.Time) > (res.IntervalAvg + res.IntervalStdDev) {
		msg := `Task "` + name + `" has not recurred in time`
		idl.Warn(msg, lastBegin)
		negativeCheck(msg)
	}

	return nil
}

func negativeCheck(nagiosMsg string) {
	nagios.SetFailed(nagiosMsg)
}
