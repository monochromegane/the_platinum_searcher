package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/terminal"
	"github.com/monochromegane/the_platinum_searcher/search"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"os"
	"runtime"
	"strings"
)

const version = "1.0.0"

var opts option.Option

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "pt"
	parser.Usage = "[OPTIONS] PATTERN [PATH]"

	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

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
		root = strings.TrimRight(args[1], "\"")
		_, err := os.Lstat(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	searcher := search.Searcher{root, args[0], &opts}
	err = searcher.Search()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
