package the_platinum_searcher

type match struct {
	path  string
	lines []line
}

type line struct {
	num  int
	text string
}

func (m *match) add(num int, text string) {
	m.lines = append(m.lines, line{num, text})
}
