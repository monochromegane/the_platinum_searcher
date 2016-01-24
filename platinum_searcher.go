package the_platinum_searcher

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/monochromegane/conflag"
	"github.com/monochromegane/go-home"
	"github.com/monochromegane/terminal"
)

const version = "2.1.0"

const (
	ExitCodeOK = iota
	ExitCodeError
)

var opts Option

type PlatinumSearcher struct {
	Out, Err io.Writer
}

func (p PlatinumSearcher) Run(args []string) int {

	parser := newOptionParser(&opts)

	conflag.LongHyphen = true
	conflag.BoolValue = false
	for _, c := range []string{filepath.Join(home.Dir(), ".ptconfig.toml"), ".ptconfig.toml"} {
		if args, err := conflag.ArgsFrom(c); err == nil {
			parser.ParseArgs(args)
		}
	}

	args, err := parser.ParseArgs(args)
	if err != nil {
		return ExitCodeError
	}

	if opts.Version {
		fmt.Printf("pt version %s\n", version)
		return ExitCodeOK
	}

	if len(args) == 0 && !opts.SearchOption.EnableFilesWithRegexp {
		parser.WriteHelp(p.Err)
		return ExitCodeError
	}

	if !terminal.IsTerminal(os.Stdout) {
		opts.OutputOption.EnableColor = false
		opts.OutputOption.EnableGroup = false
	}

	if p.givenStdin() && p.noRootPathIn(args) {
		opts.SearchOption.SearchStream = true
	}

	if opts.SearchOption.EnableFilesWithRegexp {
		args = append([]string{""}, args...)
	}

	if opts.OutputOption.Count {
		opts.OutputOption.Before = 0
		opts.OutputOption.After = 0
		opts.OutputOption.Context = 0
	}

	search := search{
		roots: p.rootsFrom(args),
		out:   p.Out,
	}
	if err = search.start(p.patternFrom(args)); err != nil {
		fmt.Fprintf(p.Err, "%s\n", err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (p PlatinumSearcher) patternFrom(args []string) string {
	return args[0]
}

func (p PlatinumSearcher) rootsFrom(args []string) []string {
	if len(args) > 1 {
		return args[1:]
	} else {
		return []string{"."}
	}
}

func (p PlatinumSearcher) givenStdin() bool {
	fi, err := os.Stdin.Stat()
	if runtime.GOOS == "windows" {
		if err == nil {
			return true
		}
	} else {
		if err != nil {
			return false
		}

		mode := fi.Mode()
		if (mode&os.ModeNamedPipe != 0) || mode.IsRegular() {
			return true
		}
	}
	return false
}

func (p PlatinumSearcher) noRootPathIn(args []string) bool {
	return len(args) == 1
}
