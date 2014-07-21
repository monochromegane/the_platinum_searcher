package the_platinum_searcher

import (
	"testing"
)

func TestFind(t *testing.T) {
	out := make(chan *GrepParams)
	find := find{out, &Option{}}
	go find.Start("files", &Pattern{Pattern: "go"})

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
	out := make(chan *GrepParams)
	find := find{out, &Option{VcsIgnore: []string{".vcsignore"}}}
	go find.Start("files/vcs", &Pattern{Pattern: "go"})

	for o := range out {
		for _, ignore := range Ignores {
			if o.Path == "files/vcs/"+ignore {
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
		out := make(chan *GrepParams)
		find := find{out, &Option{}}
		go find.Start("files/"+hidden.Root, &Pattern{Pattern: "go"})

		found := false
		for o := range out {
			if o.Path == "files/"+hidden.Expect {
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
	out := make(chan *GrepParams)
	find := find{out, &Option{Depth: 1}}
	go find.Start("files/depth", &Pattern{Pattern: "go"})

	for o := range out {
		if o.Path == "files/depth/dir_1/dir_2/file_3.txt" {
			t.Errorf("It should not find from over max depth.")
		}
	}
}

func TestFindWithFileSearchPattern(t *testing.T) {
	out := make(chan *GrepParams)
	find := find{out, &Option{}}
	pattern, _ := NewPattern("go", "match.txt", true, true, false)
	go find.Start("files/vcs/match", pattern)

	for o := range out {
		if o.Path == "files/vcs/match/ignore.txt" {
			t.Errorf("It should not contains file. %s", o.Path)
		}
	}
}

func TestFindWithStream(t *testing.T) {
	out := make(chan *GrepParams)
	find := find{out, &Option{SearchStream: true}}
	go find.Start(".", &Pattern{Pattern: "go"})

	for o := range out {
		if o.Path != "" {
			t.Errorf("It should not contains file. %s", o.Path)
		}
	}
}
