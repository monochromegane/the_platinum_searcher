package option

import (
	"testing"
)

func TestVcsIgnores(t *testing.T) {
	// When "VcsIgnore" is not specified
	expected := []string{".gitignore", ".hgignore", ".ptignore"}
	option := Option{}
	result := option.VcsIgnores()
	if !sliceEqual(expected, result) || !sliceEqual(expected, option.VcsIgnore) {
		t.Errorf("The result is invalid. [Expected: %v, Actual: %v]", expected, result)
	}

	// When "VcsIgnore" is specified
	expected = []string{".foo", ".bar", ".baz"}
	option = Option{VcsIgnore: expected}
	result = option.VcsIgnores()
	if !sliceEqual(expected, result) || !sliceEqual(expected, option.VcsIgnore) {
		t.Errorf("The result is invalid. [Expected: %v, Actual: %v]", expected, result)
	}
}

func sliceEqual(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) || cap(s1) != cap(s2) {
		return false
	}
	for i, v1 := range s1 {
		v2 := s2[i]
		if v1 != v2 {
			return false
		}
	}
	return true
}
