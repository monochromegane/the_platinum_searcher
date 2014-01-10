package main

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/the_platinum_searcher/search"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"os"
	"runtime"
)

var opts option.Option

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	args, _ := flags.Parse(&opts)

	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		opts.NoColor = true
		opts.NoGroup = true
	}

	var root = "."

	if len(args) == 0 {
		os.Exit(1)
	}
	if len(args) == 2 {
		root = args[1]
	}

	searcher := search.Searcher{root, args[0], &opts}
	searcher.Search()
}
