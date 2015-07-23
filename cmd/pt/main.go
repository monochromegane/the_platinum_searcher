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

const version = "1.7.8"

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
	opts.SkipVcsIgnore = opts.SkipVcsIgnores

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "pt"
	parser.Usage = "[OPTIONS] [PATTERN] [PATH]"

	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	pattern := ""
	if opts.FilesWithRegexp == "" && len(args) > 0 {
		pattern = args[0]
	}

	if opts.Pattern != "" {
		pattern = opts.Pattern
	}

	if opts.WordRegexp {
		opts.Regexp = true
		pattern = "\\b" + pattern + "\\b"
	}

	if pattern == "" && opts.FilesWithRegexp == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	opts.SearchStream = false
	if opts.FilesWithRegexp == "" && len(args) == 1 {
		fi, err := os.Stdin.Stat()
		if runtime.GOOS == "windows" {
			if err == nil {
				opts.SearchStream = true
				opts.NoGroup = true
			}
		} else {
			if err != nil {
				os.Exit(1)
			}

			mode := fi.Mode()
			if (mode&os.ModeNamedPipe != 0) || mode.IsRegular() {
				opts.SearchStream = true
				opts.NoGroup = true
			}
		}
	}

	var roots = []string{"."}
	if (opts.FilesWithRegexp == "" && len(args) >= 2) ||
		(opts.FilesWithRegexp != "" && len(args) > 0) {
		roots = []string{}
		paths := args[1:]
		if opts.FilesWithRegexp != "" {
			paths = args
		}
		for _, root := range paths {
			root = strings.TrimRight(root, "\"")
			_, err := os.Lstat(root)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			roots = append(roots, root)
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

	start := time.Now()

	searcher := pt.PlatinumSearcher{roots, pattern, &opts}
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

	if pt.FileMatchCount == 0 && pt.MatchCount == 0 {
		os.Exit(1)
	}
}
