package the_platinum_searcher

import (
	"bytes"
	"os"
	"testing"
)

func TestLineGrepOnlyMatch(t *testing.T) {
	opts := defaultOption()

	expect := `files/context/context.txt:4:go test
files/context/context.txt:6:go test
`

	if !assertLineGrep(opts, "files/context/context.txt", expect) {
		t.Errorf("Failed line grep (only match).")
	}
}

func TestLineGrepContext(t *testing.T) {
	opts := defaultOption()
	opts.OutputOption.Before = 2
	opts.OutputOption.After = 2

	expect := `files/context/context.txt:2-before
files/context/context.txt:3-before
files/context/context.txt:4:go test
files/context/context.txt:5-after
files/context/context.txt:6:go test
files/context/context.txt:7-after
files/context/context.txt:8-after
`

	if !assertLineGrep(opts, "files/context/context.txt", expect) {
		t.Errorf("Failed line grep (context).")
	}
}

func assertLineGrep(opts Option, path string, expect string) bool {
	buf := new(bytes.Buffer)
	printer := newPrinter(pattern{}, buf, opts)
	grep := newLineGrep(printer, opts)

	f, _ := os.Open(path)

	grep.grepEachLines(f, ASCII, func(b []byte) bool {
		return bytes.Contains(b, []byte("go"))
	}, func(b []byte) int { return 0 })

	return buf.String() == expect
}
