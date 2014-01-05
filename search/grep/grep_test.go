package grep

import (
	"github.com/monochromegane/the_platinum_searcher/search/file"
	"github.com/monochromegane/the_platinum_searcher/search/option"
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
}

func TestGrep(t *testing.T) {

	for _, g := range Asserts {
                in := make(chan *Params)
                out := make(chan *print.Params)
                grepper := Grepper{in, out, &option.Option{false, false}}

		go grepper.Grep()
		go receive(in, &Params{"../../files/" + g.path, g.pattern, g.fileType})
		o := <-out
		if o.Path != "../../files/" + g.path {
			t.Errorf("It should be equal ../../files/%s.", g.path)
		}
		if o.Matches[0].Match != g.match {
			t.Errorf("It should be equal %s", g.match)
		}
	}
}

func receive(in chan *Params, params *Params) {
	in <- params
	close(in)
}
