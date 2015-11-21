package main

import (
	"os"

	pt "github.com/monochromegane/the_platinum_searcher"
)

func main() {
	pt := pt.PlatinumSearcher{Out: os.Stdout, Err: os.Stderr}
	exitCode := pt.Run(os.Args[1:])
	os.Exit(exitCode)
}
