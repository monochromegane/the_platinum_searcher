package the_platinum_searcher

import (
	"bufio"
	"os"
	"sync"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type GrepParams struct {
	Path    string
	Encode  int
	Pattern *Pattern
}

type grep struct {
	In     chan *GrepParams
	Out    chan *PrintParams
	Option *Option
}

func Grep(in chan *GrepParams, out chan *PrintParams, option *Option) {
	grep := grep{
		In:     in,
		Out:    out,
		Option: option,
	}
	grep.ConcurrentStart()
}

var FilesSearched uint

func (g *grep) ConcurrentStart() {
	var wg sync.WaitGroup
	FilesSearched = 0
	sem := make(chan struct{}, g.Option.Proc)
	for arg := range g.In {
		sem <- struct{}{}
		wg.Add(1)
		FilesSearched++
		go func(g *grep, arg *GrepParams, sem chan struct{}) {
			defer wg.Done()
			g.Start(arg.Path, arg.Encode, arg.Pattern, sem)
		}(g, arg, sem)
	}
	wg.Wait()
	close(g.Out)
}

func (g *grep) Start(path string, encode int, pattern *Pattern, sem chan struct{}) {
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
	m := NewMatch(g.Option.Before, g.Option.After, g.Option.Column)
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

func getDecoder(encode int) transform.Transformer {
	switch encode {
	case EUCJP:
		return japanese.EUCJP.NewDecoder()
	case SHIFTJIS:
		return japanese.ShiftJIS.NewDecoder()
	}
	return nil
}

func getFileHandler(path string, opt *Option) (*os.File, error) {
	if opt.SearchStream {
		return os.Stdin, nil
	} else {
		return os.Open(path)
	}
}
