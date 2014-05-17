package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/terminal"
	"github.com/monochromegane/the_platinum_searcher/search"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"os"
	"runtime"
	"strings"
	"time"
)

const version = "1.5.2"

var opts option.Option

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {

	opts.Color = opts.SetEnableColor
	opts.NoColor = opts.SetDisableColor
	opts.EnableColor = true

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

	if len(args) == 0 && opts.FilesWithRegexp == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	opts.SearchStream = false
	if len(args) == 1 {
		if !terminal.IsTerminal(os.Stdin) {
			opts.SearchStream = true
			opts.NoGroup = true
		}
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

	opts.Proc = runtime.NumCPU()

	if !terminal.IsTerminal(os.Stdout) {
		opts.EnableColor = false
		opts.NoGroup = true
	}

	if opts.Context > 0 {
		opts.Before = opts.Context
		opts.After = opts.Context
	}

	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	start := time.Now()

	searcher := search.Searcher{root, pattern, &opts}
	err = searcher.Search()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if opts.Stats {
		elapsed := time.Since(start)
		fmt.Printf("%d Files Searched\n", grep.FilesSearched)
		fmt.Printf("%s Elapsed\n", elapsed)
	}
}
