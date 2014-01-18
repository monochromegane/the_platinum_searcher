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

var Ignores = []string{
	"match/ignore.txt",
	"ignore/ignore.txt",
	"absolute/ignore.txt",
}

func TestFindWithIgnore(t *testing.T) {
	out := make(chan *grep.Params)
	finder := Finder{out, &option.Option{VcsIgnore: []string{".vcsignore"}}}
	go finder.Find("../../files/vcs", "go")

	for o := range out {
		for _, ignore := range Ignores {
			if o.Path == "../../files/vcs/"+ignore {
				t.Errorf("It should not contains file.")
			}
		}
	}
}

func TestFindWhenSpecifiedHiddenDir(t *testing.T) {
	out := make(chan *grep.Params)
	finder := Finder{out, &option.Option{}}
	go finder.Find("../../files/.hidden", "go")

	found := false
	for o := range out {
		if o.Path == "../../files/.hidden"+"/hidden.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("It should find hidden dir when specified hidden dir.")
	}
}
