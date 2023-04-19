package main

import (
	"dfg/internal/app"
	"runtime"
)

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {

	app.Run()

}
