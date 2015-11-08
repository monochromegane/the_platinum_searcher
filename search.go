package the_platinum_searcher

type search struct {
	root    string
	pattern string
}

func (s search) start() {
	grepChan := make(chan string, 5000)
	done := make(chan struct{})

	go find{out: grepChan}.start(s.root)
	go grep{in: grepChan, done: done}.start(s.pattern)
	<-done
}
