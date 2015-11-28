package the_platinum_searcher

import "regexp"

type match struct {
	regexp *regexp.Regexp
	path   string
	lines  []line
}

type line struct {
	num  int
	text string
}

func (m *match) add(num int, text string) {
	m.lines = append(m.lines, line{num, text})
}

func (m match) size() int {
	return len(m.lines)
}
