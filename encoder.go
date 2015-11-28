package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type encoder struct {
	defaultReader  io.Reader
	eucJpReader    io.Reader
	shiftJISReader io.Reader
	opts           Option
}

func newEncoder(r io.Reader, opts Option) encoder {
	return encoder{
		defaultReader:  r,
		eucJpReader:    transform.NewReader(r, japanese.EUCJP.NewEncoder()),
		shiftJISReader: transform.NewReader(r, japanese.ShiftJIS.NewEncoder()),
		opts:           opts,
	}
}

func (e encoder) reader(encoding int) io.Reader {
	switch encoding {
	case EUCJP:
		return e.eucJpReader
	case SHIFTJIS:
		return e.shiftJISReader
	}
	return e.defaultReader
}
