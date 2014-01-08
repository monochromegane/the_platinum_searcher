package ignore

import (
        "testing"
)

func TestIgnorePatterns(t *testing.T) {

        patterns := IgnorePatterns("../../files/ignore", []string{"ignore.txt"})

        if patterns[0] != "pattern1" {
                t.Errorf("It should be equal %s", "pattern1")
        }
        if patterns[1] != "pattern2" {
                t.Errorf("It should be equal %s", "pattern2")
        }
}
