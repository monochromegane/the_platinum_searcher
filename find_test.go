package the_platinum_searcher

import (
	"regexp"
	"testing"
)

func TestFind(t *testing.T) {
	out := make(chan string)
	find := find{out, defaultOption()}
	go find.start([]string{"files"}, nil)

	testPath := makeAssertPaths(out)

	// Ensure these files were not returned
	if e := "files/.hidden/hidden.txt"; testPath(e) {
		t.Errorf("Found %s, It should not contains file under hidden directory.", e)
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
	out := make(chan string)
	opt := defaultOption()
	opt.SearchOption.Hidden = true
	find := find{out, opt}
	go find.start([]string{"files"}, nil)

	testPath := makeAssertPaths(out)

	// Enumerate found paths and ensure a couple of them are in there.
	if e := "files/.hidden/hidden.txt"; !testPath(e) {
		t.Errorf("Find failed to locate: %s", e)
	}
	if e := "files/.hidden/.hidden.txt"; !testPath(e) {
		t.Errorf("Find failed to locate: %s", e)
	}
}

func TestFindWithIgnore(t *testing.T) {
	out := make(chan string)
	opt := defaultOption()
	opt.SearchOption.VcsIgnore = []string{".vcsignore"}
	find := find{out, opt}
	go find.start([]string{"files/vcs"}, nil)

	testPath := makeAssertPaths(out)

	ignores := []string{
		"match/ignore.txt",
		"ignore/ignore.txt",
		"absolute/ignore.txt",
	}

	for _, ignore := range ignores {
		// Ensure these files were not returned
		if e := "files/vcs/" + ignore; testPath(e) {
			t.Errorf("Found %s, It should not contains ignore file.", e)
		}
	}
}

func TestFindWithDepth(t *testing.T) {
	out := make(chan string)
	opt := defaultOption()
	opt.SearchOption.Depth = 1
	find := find{out, opt}
	go find.start([]string{"files/depth"}, nil)

	testPath := makeAssertPaths(out)

	// Ensure these files were not returned
	if e := "files/depth/dir_1/dir_2/file_3.txt"; testPath(e) {
		t.Errorf("Found %s, It should not contains file from over max depth.", e)
	}
}

func TestFindWithFileSearchPattern(t *testing.T) {
	out := make(chan string)
	find := find{out, defaultOption()}
	go find.start([]string{"files/vcs/match"}, regexp.MustCompile("match.txt"))

	testPath := makeAssertPaths(out)

	// Ensure these files were not returned
	if e := "files/vcs/match/ignore.txt"; testPath(e) {
		t.Errorf("Found %s, It should not contains no match file.", e)
	}
}

func makeAssertPaths(ch chan string) func(f string) bool {
	var list []string
	for path := range ch {
		list = append(list, path)
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

func defaultOption() Option {
	return Option{
		OutputOption: &OutputOption{},
		SearchOption: &SearchOption{
			Depth: 25,
		},
	}
}
