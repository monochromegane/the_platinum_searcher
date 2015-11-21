package the_platinum_searcher

import "io"

const (
	ExitCodeOK = iota
	ExitCodeError
)

type PlatinumSearcher struct {
	Out, Err io.Writer
}

func (p PlatinumSearcher) Run(args []string) int {
	search := search{
		pattern: p.patternFrom(args),
		root:    p.rootFrom(args),
		out:     p.Out,
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
