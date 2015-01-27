package the_platinum_searcher

import (
	"strings"
)

type Match struct {
	Matched bool
	*Line
	beforeNum int
	afterNum  int
	setCol    bool
	Befores   []*Line
	Afters    []*Line
	Col       int
}

type Line struct {
	Num int
	Str string
}

func NewMatch(before, after int, setCol bool) *Match {
	return &Match{
		beforeNum: before,
		afterNum:  after,
		setCol:    setCol,
		Line:      &Line{},
	}
}

func (m *Match) newMatch() *Match {
	return &Match{
		beforeNum: m.beforeNum,
		afterNum:  m.afterNum,
		setCol:    m.setCol,
		Line:      &Line{},
	}
}

func (m *Match) LineNum() int {
	return m.Line.Num
}

func (m *Match) Match() string {
	return m.Line.Str
}

func (m *Match) setMatch(pattern *Pattern, num int, s string) {
	m.Line.Num = num
	m.Line.Str = s
	if m.setCol {
		if pattern.UseRegexp || pattern.IgnoreCase {
			m.Col = pattern.Regexp.FindStringIndex(m.Str)[0] + 1
		} else {
			m.Col = strings.Index(m.Str, pattern.Pattern) + 1
		}
	}
	m.Matched = true
}

func (m *Match) setBefore(num int, s string) {
	befores := m.Befores
	if len(m.Befores) >= m.beforeNum {
		befores = m.Befores[1:]
	}
	m.Befores = append(befores, &Line{num, s})
}

func (m *Match) setAfter(num int, s string) bool {
	if len(m.Afters) >= m.afterNum {
		return false
	}
	m.Afters = append(m.Afters, &Line{num, s})
	return true
}

func (m *Match) setUpNewMatch(pattern *Pattern, num int, s string) (*Match, bool) {
	// already match
	if m.Matched {
		newMatch := m.newMatch()
		newMatch.setMatch(pattern, num, s)
		return newMatch, true
	}
	m.setMatch(pattern, num, s)
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
			return m.setUpNewMatch(pattern, num, s)
		}
	} else if pattern.IgnoreCase {
		if strings.Contains(strings.ToUpper(s), strings.ToUpper(pattern.Pattern)) {
			return m.setUpNewMatch(pattern, num, s)
		}
	} else if strings.Contains(s, pattern.Pattern) {
		return m.setUpNewMatch(pattern, num, s)
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
