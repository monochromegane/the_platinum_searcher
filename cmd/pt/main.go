package main

import (
	"os"
	"runtime"

	pt "github.com/monochromegane/the_platinum_searcher"
)

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {
	pt := pt.PlatinumSearcher{Out: os.Stdout, Err: os.Stderr}
	exitCode := pt.Run(os.Args[1:])
	os.Exit(exitCode)
}
