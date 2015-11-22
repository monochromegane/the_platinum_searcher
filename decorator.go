package the_platinum_searcher

import (
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
	match(pattern []byte, line string) string
}

func newDecorator(option Option) decorator {
	if option.OutputOption.EnableColor {
		return color{}
	} else {
		return plain{}
	}
}

type color struct {
}

func (c color) path(path string) string {
	return ColorPath + path + ColorReset
}

func (c color) lineNumber(lineNum int) string {
	return ColorLineNumber + strconv.Itoa(lineNum) + ColorReset
}

func (c color) match(pattern []byte, line string) string {
	s := string(pattern)
	return strings.Replace(line, s, ColorMatch+s+ColorReset, -1)
}

type plain struct {
}

func (p plain) path(path string) string {
	return path
}

func (p plain) lineNumber(lineNum int) string {
	return strconv.Itoa(lineNum)
}

func (p plain) match(pattern []byte, line string) string {
	return line
}
