package analysis

import (
	"math"
	"sort"
	"time"

	idl "go.iondynamics.net/iDlogger"
)

func Avg(durations []time.Duration) time.Duration {
	if len(durations) < 1 {
		return 0
	}
	sum := time.Duration(0)
	for _, duration := range durations {
		sum = sum + duration
	}

	return sum / time.Duration(len(durations))

}

func StdDev(durations []time.Duration, precision time.Duration) time.Duration {
	if len(durations) < 1 {
		return 0
	}

	var diffSqSum float64

	avg := reducePrecision(Avg(durations), precision)

	for _, duration := range durations {
		diff := reducePrecision(duration, precision) - avg
		diffSq := diff * diff
		diffSqSum = diffSqSum + diffSq
	}

	stdDeviationApproximation := math.Sqrt(diffSqSum / float64(len(durations)-1))

	if !(stdDeviationApproximation > -1) {
		idl.Crit("overflow detected - reduce standard deviation precision")
	}

	return time.Duration(stdDeviationApproximation) * precision
}

func reducePrecision(dur time.Duration, unit time.Duration) float64 {
	switch unit {
	case time.Hour:
		return dur.Hours()
	case time.Minute:
		return dur.Minutes()
	case time.Second:
		return dur.Seconds()
	case time.Nanosecond:
		return float64(dur.Nanoseconds())
	default:
		return float64(dur)
	}
}

type ByTime []time.Time

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Before(a[j]) }

func DurationsBetween(times []time.Time) []time.Duration {
	durs := []time.Duration{}

	if len(times) < 1 {
		return durs
	}

	tbs := append([]time.Time{}, times...)
	sort.Sort(ByTime(tbs))

	last := tbs[0]

	for i := 1; i < len(tbs); i++ {
		durs = append(durs, tbs[i].Sub(last))
		last = tbs[i]
	}

	return durs
}

type ByDuration []time.Duration

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i] < a[j] }

func QuartileReduce(durations []time.Duration) (reduced []time.Duration) {
	if len(durations) < 1 {
		return
	}

	tbs := append([]time.Duration{}, durations...)
	sort.Sort(ByDuration(tbs))

	lowIndex := int(math.Floor(0.25*float64(len(tbs))) + 1)
	highIndex := int(math.Floor(0.75*float64(len(tbs))) + 1)

	if lowIndex < 0 || lowIndex >= len(tbs) {
		lowIndex = 0
	}

	if highIndex < 0 || highIndex >= len(tbs) {
		highIndex = len(tbs) - 1
	}

	lowest := tbs[lowIndex]
	hightest := tbs[highIndex]

	for _, v := range tbs {
		if v >= lowest && v <= hightest {
			reduced = append(reduced, v)
		}
	}

	return
}
