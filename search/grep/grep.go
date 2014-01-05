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

func (self *Grepper) Grep() {
	for arg := range self.In {

		fh, err := os.Open(arg.Path)
		f := bufio.NewReader(fh)
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		decoder := mahonia.NewDecoder(arg.Encode)

		m := make([]*print.Match, 0)

		var lineNum = 1
		for {
			buf, _, err = f.ReadLine()
			if err != nil {
				break
			}

			s := string(buf)
			if decoder != nil && arg.Encode != file.UTF8 && arg.Encode != file.ASCII {
				s = decoder.ConvertString(s)
			}
			if strings.Contains(s, arg.Pattern) {
				m = append(m, &print.Match{lineNum, s})
			}
			lineNum++
		}
		self.Out <- &print.Params{arg.Pattern, arg.Path, m}
		fh.Close()

	}
	close(self.Out)
}
