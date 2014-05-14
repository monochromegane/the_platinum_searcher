package grep

import (
	"bufio"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/match"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"github.com/monochromegane/the_platinum_searcher/search/print"
	"os"
	"sync"
	"launchpad.net/gommap"
	"strings"
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

var FilesSearched uint
func (self *Grepper) ConcurrentGrep() {
	var wg sync.WaitGroup
	FilesSearched = 0
	sem := make(chan bool, self.Option.Proc)
	for arg := range self.In {
		sem <- true
		wg.Add(1)
		FilesSearched++
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

func getFileHandler(path string, opt *option.Option) (*os.File, error) {
	if opt.SearchStream {
		return os.Stdin, nil
	} else {
		return os.Open(path)
	}
}

func (self *Grepper) Grep(path, encode string, pattern *pattern.Pattern, sem chan bool) {
	if self.Option.FilesWithRegexp != "" {
		self.Out <- &print.Params{pattern, path, nil}
		<-sem
		return
	}

	fh, err := getFileHandler(path, self.Option)
	if err != nil {
		panic(err)
	}

	matches := make([]*match.Match, 0)
	m := match.NewMatch(self.Option.Before, self.Option.After)
	if self.Option.SearchStream {
		var f *bufio.Reader
		if dec := getDecoder(encode); dec != nil {
			f = bufio.NewReader(transform.NewReader(fh, dec))
		} else {
			f = bufio.NewReader(fh)
		}

		var buf []byte
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
	} else {
		mmap, _ := gommap.Map(fh.Fd(), gommap.PROT_READ, gommap.MAP_PRIVATE|gommap.MAP_POPULATE)
		if (0 != len(mmap)) {
			var buf []byte
			if dec := getDecoder(encode); dec != nil {
				buf = transform.Bytes(dec, mmap)
			} else {
				buf = mmap
			}
			//buf can be nil due to transformation error
			if (buf != nil) {
				buf2 := strings.SplitN(string(buf), "\n", -1)
				var lineNum = 1
				for _, buf3 := range buf2 {
					if newMatch, ok := m.IsMatch(pattern, lineNum, buf3); ok {
						matches = append(matches, m)
						m = newMatch
					}
					lineNum++
				}
				if m.Matched {
					matches = append(matches, m)
				}

				//m.FindMatches(pattern, string(buf), &matches)
			}
		}
		//We are all done touching mmap. So it should be safe to unmap it now
		mmap.UnsafeUnmap()
	}
	self.Out <- &print.Params{pattern, path, matches}
	fh.Close()
	<-sem
}
