package the_platinum_searcher

import "testing"

func defaultOpts() *Option {

	return &Option{
		Depth:             25,
		NoGlobalGitIgnore: true,
	}
}

func mkFoundPaths(ch chan *GrepParams) func(f string) bool {
	var list []string
	for o := range ch {
		list = append(list, o.Path)
	}
	return func(m string) bool {
		for _, s := range list {
			if m == s {
				return true
			}
		}
		return false
	}
}

func TestFind(t *testing.T) {
	out := make(chan *GrepParams)
	find := find{out, defaultOpts()}
	go find.Start([]string{"files"}, &Pattern{Pattern: "go"})

	testPath := mkFoundPaths(out)

	// Ensure these files were not returned
	if e := ".hidden/hidden.txt"; testPath(e) {
		t.Errorf("Found %s, It should not contains file under hidden directory.", e)
	}
	if e := "binary/binary.bin"; testPath(e) {
		t.Errorf("%s should be text file.", e)
	}

	// Enumerate found paths and ensure a couple of them are in there.
	if e := "files/ascii.txt"; !testPath(e) {
		t.Errorf("Find failed to locate: %s", e)
	}

	if e := "files/depth/file_1.txt"; !testPath(e) {
		t.Errorf("Find failed to locate: %s", e)
	}
}

func TestFindWithHidden(t *testing.T) {
	out := make(chan *GrepParams)
	find := find{out, &Option{Hidden: true}}
	go find.Start([]string{"files"}, &Pattern{Pattern: "go"})

	testPath := mkFoundPaths(out)

	// Ensure these files were returned
	if e := ".hidden/hidden.txt"; testPath(e) {
		t.Errorf("Found %s, It should not contains file under hidden directory.", e)
	}
}

var Ignores = []string{
	"match/ignore.txt",
	"ignore/ignore.txt",
	"absolute/ignore.txt",
}

func TestFindWithIgnore(t *testing.T) {
	out := make(chan *GrepParams)
	opts := defaultOpts()
	opts.VcsIgnore = []string{".vcsignore"}
	find := find{out, opts}
	go find.Start([]string{"files/vcs"}, &Pattern{Pattern: "go"})

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
		go find.Start([]string{"files/" + hidden.Root}, &Pattern{Pattern: "go"})

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
	go find.Start([]string{"files/depth"}, &Pattern{Pattern: "go"})

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
	go find.Start([]string{"files/vcs/match"}, pattern)

	for o := range out {
		if o.Path == "files/vcs/match/ignore.txt" {
			t.Errorf("It should not contains file. %s", o.Path)
		}
	}
}

func TestFindWithStream(t *testing.T) {
	out := make(chan *GrepParams)
	find := find{out, &Option{SearchStream: true}}
	go find.Start([]string{"."}, &Pattern{Pattern: "go"})

	for o := range out {
		if o.Path != "" {
			t.Errorf("It should not contains file. %s", o.Path)
		}
	}
}
