package the_platinum_searcher

import (
	"io"
	"os"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

type PlatinumSearcher struct {
	Out, Err io.Writer
}

func (p PlatinumSearcher) Run(args []string) int {
	search := search{pattern: os.Args[1], root: os.Args[2]}
	search.start()
	return ExitCodeOK
}
