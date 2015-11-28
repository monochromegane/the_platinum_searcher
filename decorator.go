package the_platinum_searcher

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	ColorReset      = "\x1b[0m\x1b[K"
	ColorLineNumber = "\x1b[1;33m"  /* yellow with black background */
	ColorPath       = "\x1b[1;32m"  /* bold green */
	ColorMatch      = "\x1b[30;43m" /* black with yellow background */

	SeparatorColon = ":"
)

type decorator interface {
	path(path string) string
	lineNumber(lineNum int) string
	match(regexp *regexp.Regexp, line string) string
}

func newDecorator(pattern []byte, option Option) decorator {
	if option.OutputOption.EnableColor {
		return newColor(pattern)
	} else {
		return plain{}
	}
}

type color struct {
	from string
	to   string
}

func newColor(pattern []byte) color {
	color := color{}
	if pattern != nil {
		p := string(pattern)
		color.from = p
		color.to = ColorMatch + p + ColorReset
	} else {
		color.to = ColorMatch + "${1}" + ColorReset
	}
	return color
}

func (c color) path(path string) string {
	return ColorPath + path + ColorReset
}

func (c color) lineNumber(lineNum int) string {
	return ColorLineNumber + strconv.Itoa(lineNum) + ColorReset
}

func (c color) match(regexp *regexp.Regexp, line string) string {
	if regexp == nil {
		return strings.Replace(line, c.from, c.to, -1)
	} else {
		return regexp.ReplaceAllString(line, c.to)
	}
}

type plain struct {
}

func (p plain) path(path string) string {
	return path
}

func (p plain) lineNumber(lineNum int) string {
	return strconv.Itoa(lineNum)
}

func (p plain) match(regexp *regexp.Regexp, line string) string {
	return line
}
