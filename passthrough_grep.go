package the_platinum_searcher

import "sync"

type passthroughGrep struct {
	printer printer
}

func (g passthroughGrep) grep(path string, sem chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		<-sem
		wg.Done()
	}()
	match := match{path: path, lines: []line{line{}}}
	g.printer.print(match)
}
