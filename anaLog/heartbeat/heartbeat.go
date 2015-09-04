package heartbeat

import (
	"go.permanent.de/anaLog/anaLog/logpoint"
)

type Heartbeat struct {
	logpoint.LogPoint
	Subtask string
}

func Notice(hb Heartbeat) {
	//hb.Task
	//hb.RunId
}

func StillAlive(task, runid string) <-chan bool {
	ch := make(chan bool)
	go func(c chan bool) {
		//do stuff
		c <- false //debug
	}(ch)
	return ch
}

func Die(task, runid string) {

}
