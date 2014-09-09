package the_platinum_searcher

import (
	"testing"
)

func TestIgnorePatterns(t *testing.T) {

	patterns := IgnorePatterns("files/ignore", []string{"ignore.txt"}, -1)

	if !patterns[0].Match("pattern1", 0) {
		t.Errorf("It should be match %s", "pattern1")
	}
	if !patterns[0].Match("pattern2", 0) {
		t.Errorf("It should be match %s", "pattern2")
	}
}
