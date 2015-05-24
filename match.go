package the_platinum_searcher

import (
	"strings"
)

type Match struct {
	Matched bool
	*Line
	beforeNum int
	afterNum  int
	Befores   []*Line
	Afters    []*Line
}

type Line struct {
	Num int
	Col int
	Str string
}

func NewMatch(before, after int) *Match {
	return &Match{
		beforeNum: before,
		afterNum:  after,
		Line:      &Line{},
	}
}

func (m *Match) newMatch() *Match {
	return &Match{
		beforeNum: m.beforeNum,
		afterNum:  m.afterNum,
		Line:      &Line{},
	}
}

func (m *Match) LineNum() int {
	return m.Line.Num
}

func (m *Match) Match() string {
	return m.Line.Str
}

func (m *Match) setMatch(num int, col int, s string) {
	m.Line.Num = num
	m.Line.Col = col
	m.Line.Str = s
	m.Matched = true
}

func (m *Match) setBefore(num int, s string) {
	befores := m.Befores
	if len(m.Befores) >= m.beforeNum {
		befores = m.Befores[1:]
	}
	m.Befores = append(befores, &Line{num, 0, s})
}

func (m *Match) setAfter(num int, s string) bool {
	if len(m.Afters) >= m.afterNum {
		return false
	}
	m.Afters = append(m.Afters, &Line{num, 0, s})
	return true
}

func (m *Match) setUpNewMatch(num int, col int, s string) (*Match, bool) {
	// already match
	if m.Matched {
		newMatch := m.newMatch()
		newMatch.setMatch(num, col, s)
		return newMatch, true
	}
	m.setMatch(num, col, s)
	if m.afterNum == 0 {
		return m.newMatch(), true
	} else {
		// set after line
		return nil, false
	}
}

func (m *Match) IsMatch(pattern *Pattern, num int, s string) (*Match, bool) {
	if pattern.UseRegexp {
		if pattern.Regexp.MatchString(s) {
			col := strings.Index(s, pattern.Pattern)+1
			return m.setUpNewMatch(num, col, s)
		}
	} else if pattern.IgnoreCase {
		if strings.Contains(strings.ToUpper(s), strings.ToUpper(pattern.Pattern)) {
			col := strings.Index(strings.ToUpper(s), strings.ToUpper(pattern.Pattern))+1
			return m.setUpNewMatch(num, col, s)
		}
	} else if strings.Contains(s, pattern.Pattern) {
		col := strings.Index(s, pattern.Pattern)+1
		return m.setUpNewMatch(num, col, s)
	}
	if !m.Matched && m.beforeNum > 0 {
		m.setBefore(num, s)
	}
	if m.Matched && m.afterNum > 0 {
		if !m.setAfter(num, s) {
			newMatch := m.newMatch()
			if m.beforeNum > 0 {
				newMatch.setBefore(num, s)
			}
			return newMatch, true
		}
	}
	return nil, false
}

func (m *Match) FirstLineNum() int {
	if len(m.Befores) == 0 {
		return m.Line.Num
	} else {
		return m.Befores[0].Num
	}
}

func (m *Match) LastLineNum() int {
	if len(m.Afters) == 0 {
		return m.Line.Num
	} else {
		return m.Afters[len(m.Afters)-1].Num
	}
}
