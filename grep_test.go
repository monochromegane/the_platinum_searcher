package the_platinum_searcher

import (
	"os"
	"testing"
)

type GrepAssert struct {
	path, pattern string
	fileType      int
	match         string
}

var GrepAsserts = []GrepAssert{
	GrepAssert{"ascii.txt", "go", ASCII, "go test"},
	GrepAssert{"ja/euc-jp.txt", "go", EUCJP, "go テスト"},
	GrepAssert{"ja/shift_jis.txt", "go", SHIFTJIS, "go テスト"},
	GrepAssert{"ja/utf8.txt", "go", UTF8, "go テスト"},
	GrepAssert{"ja/broken_euc-jp.txt", "go", EUCJP, "go テスト"},
	GrepAssert{"ja/broken_shift_jis.txt", "go", SHIFTJIS, "go テスト"},
	GrepAssert{"ja/broken_utf8.txt", "go", UTF8, "go テスト"},
}

func TestGrep(t *testing.T) {

	for _, g := range GrepAsserts {
		in := make(chan *GrepParams)
		out := make(chan *PrintParams)
		grep := grep{in, out, &Option{Proc: 1}}

		pattern, _ := NewPattern(g.pattern, "", false, false, false)
		sem := make(chan struct{}, 1)
		sem <- struct{}{}
		go grep.Start("files/"+g.path, g.fileType, pattern, sem)
		o := <-out
		if o.Path != "files/"+g.path {
			t.Errorf("It should be equal files/%s.", g.path)
		}
		if o.Matches[0].Match() != g.match {
			t.Errorf("%s should be equal %s", g.path, g.match)
		}
	}
}

func TestGrepWithStream(t *testing.T) {
	fh, err := os.Open("files/ascii.txt")
	if err != nil {
		panic(err)
	}
	tempStdin := os.Stdin
	os.Stdin = fh
	defer func() { os.Stdin = tempStdin }()
	g := GrepAssert{"", "go", ASCII, "go test"}
	in := make(chan *GrepParams)
	out := make(chan *PrintParams)
	grep := grep{in, out, &Option{Proc: 1, SearchStream: true}}

	pattern, _ := NewPattern(g.pattern, "", false, false, false)
	sem := make(chan struct{}, 1)
	sem <- struct{}{}
	go grep.Start(g.path, g.fileType, pattern, sem)
	o := <-out
	if o.Path != g.path {
		t.Errorf("It should be equal %s.", g.path)
	}
	if o.Matches[0].Match() != g.match {
		t.Errorf("%s should be equal %s", g.path, g.match)
	}
}

func receive(in chan *GrepParams, params *GrepParams) {
	in <- params
	close(in)
}
