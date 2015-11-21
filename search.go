package the_platinum_searcher

import "io"

type search struct {
	root    string
	pattern string
	out     io.Writer
}

func (s search) start() {
	grepChan := make(chan string, 5000)
	done := make(chan struct{})

	go find{out: grepChan}.start(s.root)
	go grep{in: grepChan, done: done, printer: newPrinter(s.out)}.start(s.pattern)
	<-done
}
