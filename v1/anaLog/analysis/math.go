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
	if len(durations) < 2 {
		return 0
	}

	fullAvg := Avg(durations)

	//determine precision level dynamically
	if precision < 1 {
		switch {
		case fullAvg > time.Hour:
			precision = time.Minute
		case fullAvg > time.Minute:
			precision = time.Second
		case fullAvg > time.Second:
			precision = time.Millisecond
		case fullAvg > time.Millisecond:
			precision = time.Microsecond
		case fullAvg > time.Microsecond:
			precision = time.Nanosecond
		default:
			precision = time.Nanosecond
		}
	}

	avg := reducePrecision(fullAvg, precision)
	var diffSqSum int64
	for _, duration := range durations {
		diff := reducePrecision(duration, precision) - avg
		diffSq := diff * diff
		diffSqSum = diffSqSum + diffSq
	}

	stdDeviationApproximation := int64(math.Sqrt(float64(diffSqSum / int64(len(durations)-1))))

	if !(stdDeviationApproximation > -1) {
		if len(durations) > 5 {
			idl.Crit("overflow detected - adjust standard deviation precision")
		}
		stdDeviationApproximation = 0
	}

	return time.Duration(stdDeviationApproximation) * precision
}

func reducePrecision(dur time.Duration, unit time.Duration) int64 {
	return (int64(dur) / int64(unit))
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
