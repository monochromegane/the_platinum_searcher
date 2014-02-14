package grep

import (
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"github.com/monochromegane/the_platinum_searcher/search/print"
	"testing"
)

type Assert struct {
	path, pattern, fileType, match string
}

var Asserts = []Assert{
	Assert{"ascii.txt", "go", file.ASCII, "go test"},
	Assert{"ja/euc-jp.txt", "go", file.EUCJP, "go テスト"},
	Assert{"ja/shift_jis.txt", "go", file.SHIFTJIS, "go テスト"},
	Assert{"ja/utf8.txt", "go", file.UTF8, "go テスト"},
	Assert{"ja/broken_euc-jp.txt", "go", file.EUCJP, "go テスト"},
	Assert{"ja/broken_shift_jis.txt", "go", file.SHIFTJIS, "go テスト"},
	Assert{"ja/broken_utf8.txt", "go", file.UTF8, "go テスト"},
}

func TestGrep(t *testing.T) {

	for _, g := range Asserts {
		in := make(chan *Params)
		out := make(chan *print.Params)
		grepper := Grepper{in, out, &option.Option{Proc: 1}}

		pattern := pattern.NewPattern(g.pattern, false, false)
		sem := make(chan bool, 1)
		sem <- true
		go grepper.Grep("../../files/"+g.path, g.fileType, pattern, sem)
		o := <-out
		if o.Path != "../../files/"+g.path {
			t.Errorf("It should be equal ../../files/%s.", g.path)
		}
		if o.Matches[0].Match != g.match {
			t.Errorf("%s should be equal %s", g.path, g.match)
		}
	}
}

func receive(in chan *Params, params *Params) {
	in <- params
	close(in)
}
