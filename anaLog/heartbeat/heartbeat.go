package heartbeat

import (
	"go.iondynamics.net/iDhelper/randGen"

	"go.permanent.de/anaLog/anaLog/logpoint"
)

type Heartbeat struct {
	logpoint.LogPoint
	Subtask string
}

func (hb Heartbeat) UID() string {
	return hb.LogPoint.Task + "_" + hb.LogPoint.RunId + "_" + hb.Subtask + "_" + randGen.String(16)
}

var reg = NewRegistry()

func Create(hb Heartbeat) {
	reg.Create(hb)
}

func Ping(hb Heartbeat) {
	reg.Ping(hb)
}

func Wait(hb Heartbeat) {
	<-reg.GetCh(hb)
	reg.Clear(hb)
}

func Exit(hb Heartbeat) {
	reg.Clear(hb)
}
