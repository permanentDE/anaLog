package heartbeat

import (
	//"go.permanent.de/anaLog/anaLog/analysis"
	"sync"
	"time"
)

type Registry struct {
	mu  sync.RWMutex
	reg map[string][]Heartbeat
}

func NewRegistry() *Registry {
	return &Registry{reg: make(map[string][]Heartbeat)}
}

func (r *Registry) Ping(hb Heartbeat) {
	r.mu.Lock()
	defer r.mu.Unlock()

	slice := r.reg[hb.Task]
	slice = append(slice, hb)
	r.reg[hb.Task] = slice
}

func (r *Registry) SubtaskTimes(hb Heartbeat) map[string][]time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()

	subtaskTimes := make(map[string][]time.Time)

	for _, h := range r.reg[hb.Task] {
		times := subtaskTimes[h.Subtask]
		times = append(times, h.Time)
		subtaskTimes[h.Subtask] = times
	}

	return subtaskTimes
}

func (r *Registry) Clear(hb Heartbeat) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.reg, hb.Task)
}
