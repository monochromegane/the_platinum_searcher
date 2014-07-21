package the_platinum_searcher

import (
	"testing"
)

type PatternAssert struct {
	Pattern string
	Expect  bool
}

func TestIgnoreCaseWithSmartCase(t *testing.T) {

	asserts := []PatternAssert{
		PatternAssert{"lowercase", true},
		PatternAssert{"Uppercase", false},
	}

	for _, assert := range asserts {
		pattern, _ := NewPattern(assert.Pattern, "", true, true, false)
		if pattern.IgnoreCase != assert.Expect {
			t.Errorf("When pattern is %s, ignore case should be %t.", assert.Pattern, assert.Expect)
		}
	}

}
