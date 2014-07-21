package the_platinum_searcher

import (
	"bufio"
	"os"
	"sync"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
)

type GrepParams struct {
	Path, Encode string
	Pattern      *pattern.Pattern
}

type Grepper struct {
	In     chan *GrepParams
	Out    chan *PrintParams
	Option *option.Option
}

var FilesSearched uint

func (g *Grepper) ConcurrentGrep() {
	var wg sync.WaitGroup
	FilesSearched = 0
	sem := make(chan bool, g.Option.Proc)
	for arg := range g.In {
		sem <- true
		wg.Add(1)
		FilesSearched++
		go func(g *Grepper, arg *GrepParams, sem chan bool) {
			defer wg.Done()
			g.Grep(arg.Path, arg.Encode, arg.Pattern, sem)
		}(g, arg, sem)
	}
	wg.Wait()
	close(g.Out)
}

func getDecoder(encode string) transform.Transformer {
	switch encode {
	case EUCJP:
		return japanese.EUCJP.NewDecoder()
	case SHIFTJIS:
		return japanese.ShiftJIS.NewDecoder()
	}
	return nil
}

func getFileHandler(path string, opt *option.Option) (*os.File, error) {
	if opt.SearchStream {
		return os.Stdin, nil
	} else {
		return os.Open(path)
	}
}

func (g *Grepper) Grep(path, encode string, pattern *pattern.Pattern, sem chan bool) {
	if g.Option.FilesWithRegexp != "" {
		g.Out <- &PrintParams{pattern, path, nil}
		<-sem
		return
	}

	fh, err := getFileHandler(path, g.Option)
	if err != nil {
		panic(err)
	}

	var f *bufio.Reader
	if dec := getDecoder(encode); dec != nil {
		f = bufio.NewReader(transform.NewReader(fh, dec))
	} else {
		f = bufio.NewReader(fh)
	}

	var buf []byte
	matches := make([]*Match, 0)
	m := NewMatch(g.Option.Before, g.Option.After)
	var lineNum = 1
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			break
		}
		if newMatch, ok := m.IsMatch(pattern, lineNum, string(buf)); ok {
			matches = append(matches, m)
			m = newMatch
		}
		lineNum++
	}
	if m.Matched {
		matches = append(matches, m)
	}
	g.Out <- &PrintParams{pattern, path, matches}
	fh.Close()
	<-sem
}
