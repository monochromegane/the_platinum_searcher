package grep

import (
	"bufio"
	"code.google.com/p/mahonia"
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/print"
	"os"
	"strings"
)

type Params struct {
	Path, Pattern, Encode string
}

type Grepper struct {
	In     chan *Params
	Out    chan *print.Params
	Option *option.Option
}

func (self *Grepper) ConcurrentGrep() {
	sem := make(chan bool, self.Option.Proc)
	for arg := range self.In {
		sem <- true
		go self.Grep(arg.Path, arg.Encode, arg.Pattern, sem)
	}
	for {
		if len(sem) == 0 {
			break
		}
	}
	close(self.Out)
}

func (self *Grepper) Grep(path, encode, pattern string, finish chan bool) {

	fh, err := os.Open(path)
	f := bufio.NewReader(fh)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	decoder := mahonia.NewDecoder(encode)

	m := make([]*print.Match, 0)

	var lineNum = 1
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			break
		}

		s := string(buf)
		if decoder != nil && encode != file.UTF8 && encode != file.ASCII {
			s = decoder.ConvertString(s)
		}
		if strings.Contains(s, pattern) {
			m = append(m, &print.Match{lineNum, s})
		}
		lineNum++
	}
	self.Out <- &print.Params{pattern, path, m}
	fh.Close()
	<-finish
}
