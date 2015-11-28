package the_platinum_searcher

type match struct {
	path  string
	lines []line
}

type line struct {
	num     int
	text    string
	matched bool
}

func (m *match) add(num int, text string, matched bool) {
	m.lines = append(m.lines, line{
		num:     num,
		text:    text,
		matched: matched,
	})
}

func (m match) size() int {
	return len(m.lines)
}
