package the_platinum_searcher

import (
	"bytes"
	"strings"
	"testing"
)

func TestFixedGrep(t *testing.T) {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.EnableGroup = false

	pattern, _ := newPattern("go", opts)

	paths := []string{
		"ascii.txt",
		"ja/utf8.txt",
		"ja/euc-jp.txt",
		"ja/shift_jis.txt",
		"ja/broken_utf8.txt",
		"ja/broken_euc-jp.txt",
		"ja/broken_shift_jis.txt",
	}

	asserts := []string{
		"ascii.txt:2:go test",
		"ja/utf8.txt:2:go テスト",
		"ja/euc-jp.txt:2:go テスト",
		"ja/shift_jis.txt:2:go テスト",
		"ja/broken_utf8.txt:2:go テスト",
		"ja/broken_euc-jp.txt:2:go テスト",
		"ja/broken_shift_jis.txt:2:go テスト",
	}

	if !assertGrep(pattern, opts, paths, asserts) {
		t.Errorf("Grep result should contain assserts.")
	}

}

func TestExtendedGrep(t *testing.T) {
	opts := defaultOption()
	opts.OutputOption.EnableColor = false
	opts.OutputOption.EnableGroup = false
	opts.SearchOption.Regexp = true

	pattern, _ := newPattern("g.*", opts)

	paths := []string{
		"ascii.txt",
		"ja/utf8.txt",
		"ja/euc-jp.txt",
		"ja/shift_jis.txt",
		"ja/broken_utf8.txt",
		"ja/broken_euc-jp.txt",
		"ja/broken_shift_jis.txt",
	}

	asserts := []string{
		"ascii.txt:2:go test",
		"ja/utf8.txt:2:go テスト",
		"ja/euc-jp.txt:2:go テスト",
		"ja/shift_jis.txt:2:go テスト",
		"ja/broken_utf8.txt:2:go テスト",
		"ja/broken_euc-jp.txt:2:go テスト",
		"ja/broken_shift_jis.txt:2:go テスト",
	}

	if !assertGrep(pattern, opts, paths, asserts) {
		t.Errorf("Grep result should contain assserts.")
	}

}

func assertGrep(pattern pattern, opts Option, paths, asserts []string) bool {
	buf := new(bytes.Buffer)
	printer := newPrinter(pattern, buf, opts)

	in := make(chan string)
	done := make(chan struct{})
	grep := newGrep(pattern, in, done, opts, printer)
	go grep.start()

	for _, path := range paths {
		in <- "files/" + path
	}
	close(in)
	<-done

	result := buf.String()
	for _, assert := range asserts {
		if !strings.Contains(result, assert) {
			return false
		}
	}
	return true
}
