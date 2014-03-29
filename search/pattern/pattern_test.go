package pattern

import (
	"testing"
)

type Assert struct {
	Pattern string
	Expect  bool
}

func TestIgnoreCaseWithSmartCase(t *testing.T) {

	asserts := []Assert{
		Assert{"lowercase", true},
		Assert{"Uppercase", false},
	}

	for _, assert := range asserts {
		pattern, _ := NewPattern(assert.Pattern, "", true, true)
		if pattern.IgnoreCase != assert.Expect {
			t.Errorf("When pattern is %s, ignore case should be %t.", assert.Pattern, assert.Expect)
		}
	}

}
