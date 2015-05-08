package main

import (
	anaLog "permanent/anaLog/v1"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	anaLog.Run()
}
