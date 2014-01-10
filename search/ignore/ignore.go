package ignore

import (
	"bufio"
	"os"
	"strings"
)

type Ignore struct {
	Patterns, Matches []string
}

func IgnorePatterns(path string, ignores []string) []string {
	var patterns []string
	for _, ignore := range ignores {
		file, err := os.Open(path + string(os.PathSeparator) + ignore)
		if err != nil {
			continue
		}
		reader := bufio.NewReader(file)
		buf := make([]byte, 1024)

		for {
			buf, _, err = reader.ReadLine()
			if err != nil {
				break
			}
			s := strings.Trim(string(buf), " ")

			if len(s) == 0 || strings.HasPrefix(s, "#") {
				continue
			}
			patterns = append(patterns, s)
		}
		file.Close()
	}
	return patterns
}
