package randGen

import (
	"crypto/rand"
)

var Source = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Bytes(n int) []byte {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = Source[b%byte(len(Source))]
	}
	return bytes
}

func String(n int) string {
	return string(Bytes(n))
}
