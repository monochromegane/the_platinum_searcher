package main

import (
	"github.com/monochromegane/the_platinum_searcher/searcher"
	"os"
)

func main() {

	var root = "."

	if len(os.Args) == 1 {
		os.Exit(1)
	}
	if len(os.Args) == 3 {
		root = os.Args[2]
	}

	searcher := pt.Searcher{root, os.Args[1]}
	searcher.Search()
}
