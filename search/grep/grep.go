package grep

import (
	"bufio"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"github.com/monochromegane/the_platinum_searcher/search/print"
	"os"
	"strings"
	"sync"
)

type Params struct {
	Path, Encode string
	Pattern      *pattern.Pattern
}

type Grepper struct {
	In     chan *Params
	Out    chan *print.Params
	Option *option.Option
}

func (self *Grepper) ConcurrentGrep() {
	var wg sync.WaitGroup
	sem := make(chan bool, self.Option.Proc)
	for arg := range self.In {
		sem <- true
		wg.Add(1)
		go func(self *Grepper, arg *Params, sem chan bool) {
			defer wg.Done()
			self.Grep(arg.Path, arg.Encode, arg.Pattern, sem)
		}(self, arg, sem)
	}
	wg.Wait()
	close(self.Out)
}

func getDecoder(encode string) transform.Transformer {
	switch encode {
	case file.EUCJP:
		return japanese.EUCJP.NewDecoder()
	case file.SHIFTJIS:
		return japanese.ShiftJIS.NewDecoder()
	}
	return nil
}

func (self *Grepper) Grep(path, encode string, pattern *pattern.Pattern, sem chan bool) {
	fh, err := os.Open(path)
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
	m := make([]*print.Match, 0)
	var lineNum = 1
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			break
		}

		s := string(buf)
		if pattern.IgnoreCase {
			if pattern.Regexp.MatchString(s) {
				m = append(m, &print.Match{lineNum, s})
			}
		} else if strings.Contains(s, pattern.Pattern) {
			m = append(m, &print.Match{lineNum, s})
		}
		lineNum++
	}
	self.Out <- &print.Params{pattern, path, m}
	fh.Close()
	<-sem
}
