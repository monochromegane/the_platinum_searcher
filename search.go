package the_platinum_searcher

import (
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"github.com/monochromegane/the_platinum_searcher/search/print"
)

type Searcher struct {
	Root, Pattern string
	Option        *option.Option
}

func (s *Searcher) Search() error {
	pattern, err := s.pattern()
	if err != nil {
		return err
	}
	grep := make(chan *GrepParams, s.Option.Proc)
	match := make(chan *print.Params, s.Option.Proc)
	done := make(chan bool)
	go s.find(grep, pattern)
	go s.grep(grep, match)
	go s.print(match, done)
	<-done
	return nil
}

func (s *Searcher) pattern() (*pattern.Pattern, error) {
	fileRegexp := s.Option.FileSearchRegexp
	if s.Option.FilesWithRegexp != "" {
		fileRegexp = s.Option.FilesWithRegexp
	}
	return pattern.NewPattern(
		s.Pattern,
		fileRegexp,
		s.Option.SmartCase,
		s.Option.IgnoreCase,
		s.Option.Regexp,
	)
}

func (s *Searcher) find(out chan *GrepParams, pattern *pattern.Pattern) {
	finder := Finder{out, s.Option}
	finder.Find(s.Root, pattern)
}

func (s *Searcher) grep(in chan *GrepParams, out chan *print.Params) {
	grepper := Grepper{in, out, s.Option}
	grepper.ConcurrentGrep()
}

func (s *Searcher) print(in chan *print.Params, done chan bool) {
	printer := print.NewPrinter(in, done, s.Option)
	printer.Print()
}
