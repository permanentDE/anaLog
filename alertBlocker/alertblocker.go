package alertBlocker

import (
	"sync"
)

var mu sync.RWMutex
var reg = make(map[string][]string)

//IsUnknown returns true if the given task currently doesn't have the given problem
func IsUnknown(task, problem string) bool {
	mu.RLock()
	defer mu.RUnlock()

	slice, ok := reg[task]
	if ok {
		for _, val := range slice {
			if val == problem {
				return false
			}
		}
	}
	return true
}

//Occurred should be called when the given task has the given problem
func Occurred(task, problem string) {
	mu.Lock()
	defer mu.Unlock()

	slice := reg[task]
	slice = append(slice, problem)
	reg[task] = slice
}

//Resolved should be called when the given task does not have the given problem anymore
func Resolved(task, problem string) {
	mu.Lock()
	defer mu.Unlock()

	newSlice := []string{}
	slice := reg[task]
	for _, val := range slice {
		if val != problem {
			newSlice = append(newSlice, val)
		}
	}

	if len(newSlice) > 0 {
		reg[task] = newSlice
	} else {
		delete(reg, task)
	}
}

//AllProblems returns a map with all known problems. Key = Taskname ; Value = Slice of problems
func AllProblems() map[string][]string {
	mu.RLock()
	defer mu.RUnlock()

	m := make(map[string][]string)
	for k, v := range reg {
		m[k] = v
	}
	return m
}
