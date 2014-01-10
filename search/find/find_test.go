package find

import (
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"testing"
)

func TestFind(t *testing.T) {
	out := make(chan *grep.Params)
	finder := Finder{out, &option.Option{}}
	go finder.Find("../../files", "go")

	for o := range out {
		if o.Path == ".hidden/hidden.txt" {
			t.Errorf("It should not contains file under hidden directory.")
		}
		if o.Path == "binary/binary.bin" {
			t.Errorf("It should be text file.")
		}
	}

}
