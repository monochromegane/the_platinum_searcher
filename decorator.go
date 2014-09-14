package the_platinum_searcher

import (
	"fmt"
	"strings"
)

const (
	ColorReset      = "\x1b[0m\x1b[K"
	ColorLineNumber = "\x1b[1;33m"  /* yellow with black background */
	ColorPath       = "\x1b[1;32m"  /* bold green */
	ColorMatch      = "\x1b[30;43m" /* black with yellow background */
)

type Decorator interface {
	path(path string) string
	lineNumber(lineNum int, sep string) string
	match(pattern *Pattern, line *Line) string
}

func newDecorator(option *Option) Decorator {
	if option.EnableColor {
		return Color{}
	} else {
		return Plain{}
	}
}

type Color struct {
}

func (c Color) path(path string) string {
	return fmt.Sprintf("%s%s%s", ColorPath, path, ColorReset)
}

func (c Color) lineNumber(lineNum int, sep string) string {
	return fmt.Sprintf("%s%d%s%s", ColorLineNumber, lineNum, ColorReset, sep)
}

func (c Color) match(pattern *Pattern, line *Line) string {
	if pattern.UseRegexp || pattern.IgnoreCase {
		return pattern.Regexp.ReplaceAllString(line.Str, ColorMatch+"${1}"+ColorReset)
	} else {
		return strings.Replace(line.Str, pattern.Pattern, ColorMatch+pattern.Pattern+ColorReset, -1)
	}
}

type Plain struct {
}

func (p Plain) path(path string) string {
	return path
}

func (p Plain) lineNumber(lineNum int, sep string) string {
	return fmt.Sprintf("%d%s", lineNum, sep)
}

func (p Plain) match(pattern *Pattern, line *Line) string {
	return line.Str
}
