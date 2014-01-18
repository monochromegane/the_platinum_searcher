package search

import (
	"github.com/monochromegane/the_platinum_searcher/search/find"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/print"
)

type Searcher struct {
	Root, Pattern string
	Option        *option.Option
}

func (self *Searcher) Search() {
	grep := make(chan *grep.Params, self.Option.Proc)
	match := make(chan *print.Params, self.Option.Proc)
	done := make(chan bool)
	go self.find(grep)
	go self.grep(grep, match)
	go self.print(match, done)
	<-done
}

func (self *Searcher) find(out chan *grep.Params) {
	finder := find.Finder{out, self.Option}
	finder.Find(self.Root, self.Pattern)
}

func (self *Searcher) grep(in chan *grep.Params, out chan *print.Params) {
	grepper := grep.Grepper{in, out, self.Option}
	grepper.ConcurrentGrep()
}

func (self *Searcher) print(in chan *print.Params, done chan bool) {
	printer := print.Printer{in, done, self.Option}
	printer.Print()
}
