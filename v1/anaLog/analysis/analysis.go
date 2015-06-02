package analysis

import (
	"math"
	"sort"
	"time"

	idl "go.iondynamics.net/iDlogger"
)

func Avg(durations []time.Duration) time.Duration {
	sum := time.Duration(0)
	for _, duration := range durations {
		sum = sum + duration
	}

	return sum / time.Duration(len(durations))

}

func StdDev(durations []time.Duration, precision time.Duration) time.Duration {
	var diffSqSum float64

	avg := Avg(durations)

	for _, duration := range durations {
		diff := (duration - avg) / precision
		diffSq := diff * diff
		diffSqSum = diffSqSum + float64(diffSq)
	}

	stdDeviationApproximation := math.Sqrt(diffSqSum / float64(len(durations)-1))

	if !(stdDeviationApproximation > 0) {
		idl.Crit("overflow detected - reduce standard deviation precision")
	}

	return time.Duration(int(stdDeviationApproximation)) * precision
}

type ByDuration []time.Duration

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i] < a[j] }

func QuartileReduce(durations []time.Duration) (reduced []time.Duration) {
	tbs := durations
	sort.Sort(ByDuration(tbs))

	lowest := tbs[int(math.Floor(0.25*float64(len(durations)))+1)]
	hightest := tbs[int(math.Floor(0.75*float64(len(durations)))+1)]

	for _, v := range tbs {
		if v >= lowest && v <= hightest {
			reduced = append(reduced, v)
		}
	}

	return
}
