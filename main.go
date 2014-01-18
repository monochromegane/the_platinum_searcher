package main

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/terminal"
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

	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	parser.Name = "pt"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"

	opts.Proc = runtime.NumCPU()

	if !terminal.IsTerminal(os.Stdout) {
		opts.NoColor = true
		opts.NoGroup = true
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	var root = "."
	if len(args) == 2 {
		root = args[1]
	}

	searcher := search.Searcher{root, args[0], &opts}
	searcher.Search()
}
