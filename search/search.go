package search

import (
	"github.com/monochromegane/the_platinum_searcher/search/find"
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"github.com/monochromegane/the_platinum_searcher/search/print"
)

type Searcher struct {
	Root, Pattern string
	Option        *option.Option
}

func (self *Searcher) Search() error {
	pattern, err:= self.pattern()
	if err != nil {
		return err
	}
	grep := make(chan *grep.Params, self.Option.Proc)
	match := make(chan *print.Params, self.Option.Proc)
	done := make(chan bool)
	go self.find(grep, pattern)
	go self.grep(grep, match)
	go self.print(match, done)
	<-done
	return nil
}

func (self *Searcher) pattern() (*pattern.Pattern, error) {
	return pattern.NewPattern(
		self.Pattern,
		self.Option.FileSearchRegexp,
		self.Option.SmartCase,
		self.Option.IgnoreCase,
	)
}

func (self *Searcher) find(out chan *grep.Params, pattern *pattern.Pattern) {
	finder := find.Finder{out, self.Option}
	finder.Find(self.Root, pattern)
}

func (self *Searcher) grep(in chan *grep.Params, out chan *print.Params) {
	grepper := grep.Grepper{in, out, self.Option}
	grepper.ConcurrentGrep()
}

func (self *Searcher) print(in chan *print.Params, done chan bool) {
	printer := print.Printer{in, done, self.Option}
	printer.Print()
}
