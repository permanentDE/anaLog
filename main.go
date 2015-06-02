package main

import (
	anaLog "go.permanent.de/anaLog/v1"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	anaLog.Run()
}
