package main

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/the_platinum_searcher/searcher"
	"os"
	"runtime"
)

var opts pt.SearchOption

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	args, _ := flags.Parse(&opts)

	var root = "."

	if len(args) == 0 {
		os.Exit(1)
	}
	if len(args) == 2 {
		root = args[1]
	}

	searcher := pt.Searcher{root, args[0], &opts}
	searcher.Search()
}
