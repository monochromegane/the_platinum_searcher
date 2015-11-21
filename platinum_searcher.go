package the_platinum_searcher

import (
	"fmt"
	"io"
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

	search := search{
		pattern: p.patternFrom(args),
		root:    p.rootFrom(args),
		out:     p.Out,
		opts:    opts,
	}
	search.start()
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
