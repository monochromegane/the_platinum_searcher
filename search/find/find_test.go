package find

import (
	"github.com/monochromegane/the_platinum_searcher/search/grep"
	"github.com/monochromegane/the_platinum_searcher/search/option"
	"github.com/monochromegane/the_platinum_searcher/search/pattern"
	"testing"
)

func TestFind(t *testing.T) {
	out := make(chan *grep.Params)
	finder := Finder{out, &option.Option{}}
	go finder.Find("../../files", &pattern.Pattern{Pattern: "go"})

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
	go finder.Find("../../files/vcs", &pattern.Pattern{Pattern: "go"})

	for o := range out {
		for _, ignore := range Ignores {
			if o.Path == "../../files/vcs/"+ignore {
				t.Errorf("It should not contains file.")
			}
		}
	}
}

type Hidden struct {
	Root, Expect string
}

var Hiddens = []Hidden{
	Hidden{".hidden", ".hidden/hidden.txt"},
	Hidden{".hidden/.hidden.txt", ".hidden/.hidden.txt"},
}

func TestFindWhenSpecifiedHiddenFile(t *testing.T) {
	for _, hidden := range Hiddens {
		out := make(chan *grep.Params)
		finder := Finder{out, &option.Option{}}
		go finder.Find("../../files/"+hidden.Root, &pattern.Pattern{Pattern: "go"})

		found := false
		for o := range out {
			if o.Path == "../../files/"+hidden.Expect {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("It should find hidden dir or file when specified hidden one.")
		}
	}
}

func TestFindWithDepth(t *testing.T) {
	out := make(chan *grep.Params)
	finder := Finder{out, &option.Option{Depth: 1}}
	go finder.Find("../../files/depth", &pattern.Pattern{Pattern: "go"})

	for o := range out {
		if o.Path == "../../files/depth/dir_1/dir_2/file_3.txt" {
			t.Errorf("It should not find from over max depth.")
		}
	}
}
