package mongo

import (
	"time"
)

type InterfaceTransportObject struct {
	If   interface{}
	Time time.Time
}

func InterfaceTransport(obj interface{}) InterfaceTransportObject {
	return InterfaceTransportObject{
		If:   obj,
		Time: time.Now(),
	}
}

type BytesTransportObject struct {
	Bytes []byte
	Time  time.Time
}

func BytesTransport(b []byte) BytesTransportObject {
	return BytesTransportObject{
		Bytes: b,
		Time:  time.Now(),
	}
}
