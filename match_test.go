package the_platinum_searcher

import (
	"testing"

	"github.com/monochromegane/the_platinum_searcher/search/pattern"
)

func testIsMatch(t *testing.T) {

	match := NewMatch(0, 0)

	// not use regexp
	p, _ := pattern.NewPattern("go", "", false, false, false)
	lines := []string{"go", "GO", "Go", "oo"}
	for index, line := range lines {
		_, ok := match.IsMatch(p, index+1, line)
		if ok {
			if match.Str != "go" {
				t.Errorf("It should be equal %s, but %s.", "go", match.Str)
			}
		}
	}

	// ignore case
	p, _ = pattern.NewPattern("go", "", false, true, false)
	lines = []string{"go", "GO", "Go", "oo"}
	for index, line := range lines {
		_, ok := match.IsMatch(p, index+1, line)
		if ok {
			if match.Str != "go" || match.Str != "GO" || match.Str != "Go" {
				t.Errorf("It should be equal %s, but %s.", "go|Go", match.Str)
			}
		}
	}

	// use regexp
	p, _ = pattern.NewPattern("go|Go", "", false, false, true)
	lines = []string{"go", "GO", "Go", "oo"}
	for index, line := range lines {
		_, ok := match.IsMatch(p, index+1, line)
		if ok {
			if match.Str != "go" || match.Str != "Go" {
				t.Errorf("It should be equal %s, but %s.", "go|Go", match.Str)
			}
		}
	}

}

func TestMatch(t *testing.T) {

	pattern, _ := pattern.NewPattern("go", "", false, false, false)
	match := NewMatch(1, 1)

	lines := []string{
		"before",
		"go match",
		"after",
	}

	for index, line := range lines {
		_, ok := match.IsMatch(pattern, index+1, line)
		if ok {
			if match.Str != "go match" {
				t.Errorf("It should be equal %s, but %s.", "go match", match.Str)
			}
			if match.Befores[0].Str != "before" {
				t.Errorf("It should be equal %s, but %s.", "before", match.Befores[0].Str)
			}
			if match.Afters[0].Str != "after" {
				t.Errorf("It should be equal %s, but %s.", "after", match.Afters[0].Str)
			}
		}
	}
}

func TestMatchWhenContextAndMatchDuplicate(t *testing.T) {

	pattern, _ := pattern.NewPattern("go", "", false, false, false)
	match := NewMatch(1, 1)

	lines := []string{
		"before",
		"go match 1",
		"go match 2",
		"after",
	}

	for index, line := range lines {
		newMatch, ok := match.IsMatch(pattern, index+1, line)
		if ok && match.Str == "go match 1" {
			if len(match.Befores) != 1 {
				t.Errorf("It should be equal %d, but %d.", 1, len(match.Befores))
			}
			if len(match.Afters) != 0 {
				t.Errorf("It should be equal %d, but %d.", 0, len(match.Afters))
			}
			if newMatch.Matched != true {
				t.Errorf("It should be equal %b, but %b.", true, match.Matched)
			}
		}
		if ok && match.Str == "go match 2" {
			if len(match.Befores) != 0 {
				t.Errorf("It should be equal %d, but %d.", 0, len(match.Befores))
			}
			if len(match.Afters) != 1 {
				t.Errorf("It should be equal %d, but %d.", 1, len(match.Afters))
			}
			if newMatch.Matched != false {
				t.Errorf("It should be equal %b, but %b.", false, match.Matched)
			}
		}

	}
}
