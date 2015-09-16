package heartbeat

import (
	"sort"
	"sync"
	"time"

	"go.permanent.de/anaLog/analysis"
)

type Registry struct {
	mu  sync.RWMutex
	hbs map[string][]Heartbeat
	chs map[string]chan bool
}

func NewRegistry() *Registry {
	return &Registry{hbs: make(map[string][]Heartbeat), chs: make(map[string]chan bool)}
}

func (r *Registry) Create(hb Heartbeat) chan bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := make(chan bool)
	r.chs[hb.Task+hb.RunId] = ch
	r.hbs[hb.Task+hb.RunId] = []Heartbeat{hb}

	return ch
}

func (r *Registry) Ping(hb Heartbeat) {
	if !r.Exists(hb) {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	slice := r.hbs[hb.Task+hb.RunId]
	slice = append(slice, hb)
	r.hbs[hb.Task+hb.RunId] = slice
}

func (r *Registry) SubtaskTimes(hb Heartbeat) map[string][]time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()

	subtaskTimes := make(map[string][]time.Time)

	for _, h := range r.hbs[hb.Task+hb.RunId] {
		times := subtaskTimes[h.Subtask]
		times = append(times, h.Time)
		subtaskTimes[h.Subtask] = times
	}

	return subtaskTimes
}

func (r *Registry) Clear(hb Heartbeat) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.hbs, hb.Task+hb.RunId)
}

func (r *Registry) GetCh(hb Heartbeat) chan bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ch, ok := r.chs[hb.Task+hb.RunId]
	if ok {
		go r.notifier(ch, hb)
	} else {
		ch = make(chan bool)
		go func(c chan bool) { c <- true }(ch)
	}

	return ch
}

func (r *Registry) Exists(hb Heartbeat) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.hbs[hb.Task+hb.RunId]
	return ok
}

func (r *Registry) notifier(ch chan bool, hb Heartbeat) {
	for {
		sleep, exit := func(c chan bool, h Heartbeat, reg *Registry) (time.Duration, bool) {
			var wait4 time.Duration
			exit := true

			for _, times := range reg.SubtaskTimes(h) {
				if len(times) > 1 {
					durs := analysis.DurationsBetween(times)
					maxDiff := analysis.Avg(durs) + analysis.StdDev(durs, 0)
					sort.Sort(analysis.ByTime(times))
					if time.Now().Sub(times[len(times)-1]) < maxDiff {
						exit = false
						if wait4 < maxDiff {
							wait4 = maxDiff
						}
					}
				}
			}

			if exit {
				c <- true
			}
			return wait4, exit
		}(ch, hb, r)
		if exit == true {
			break
		} else {

		}
		<-time.After(sleep)
	}
}
