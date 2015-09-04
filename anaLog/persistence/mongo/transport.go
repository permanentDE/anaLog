package mongo

import (
	"time"
)

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
