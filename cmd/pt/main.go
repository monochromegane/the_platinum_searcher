package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/monochromegane/terminal"
	pt "github.com/monochromegane/the_platinum_searcher"
)

const version = "1.7.1"

var opts pt.Option

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
		if !opts.ForceColor {
			opts.EnableColor = false
		}
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

	if opts.WordRegexp {
		opts.Regexp = true
		pattern = "\\b" + pattern + "\\b"
	}

	start := time.Now()

	searcher := pt.PlatinumSearcher{root, pattern, &opts}
	err = searcher.Search()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if opts.Stats {
		elapsed := time.Since(start)
		fmt.Printf("%d Files Searched\n", pt.FilesSearched)
		fmt.Printf("%s Elapsed\n", elapsed)
	}

	if pt.FileMatchCount == 0 {
		os.Exit(1)
	}
}
