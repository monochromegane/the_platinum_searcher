package the_platinum_searcher

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/monochromegane/conflag"
	"github.com/monochromegane/go-home"
	"github.com/monochromegane/terminal"
)

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
		fmt.Fprintf(p.Err, "%s\n", err)
		return ExitCodeError
	}

	if opts.Version {
		fmt.Printf("pt version %s\n", "2.0.0")
		return ExitCodeOK
	}

	if len(args) == 0 {
		parser.WriteHelp(p.Err)
		return ExitCodeError
	}

	if !terminal.IsTerminal(os.Stdout) {
		opts.OutputOption.EnableColor = false
		opts.OutputOption.EnableGroup = false
	}

	search := search{
		root: p.rootFrom(args),
		out:  p.Out,
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

func (p PlatinumSearcher) rootFrom(args []string) string {
	var root string
	if len(args) > 1 {
		root = args[1]
	} else {
		root = "."
	}
	return root
}
