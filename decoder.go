package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type decoder struct {
	defaultWriter  io.Writer
	eucJpWriter    io.Writer
	shiftJISWriter io.Writer
	opts           Option
}

func newDecoder(w io.Writer, opts Option) decoder {
	return decoder{
		defaultWriter:  w,
		eucJpWriter:    transform.NewWriter(w, japanese.EUCJP.NewDecoder()),
		shiftJISWriter: transform.NewWriter(w, japanese.ShiftJIS.NewDecoder()),
	}
}

func (d decoder) writer(encoding int) io.Writer {
	switch encoding {
	case EUCJP:
		return d.eucJpWriter
	case SHIFTJIS:
		return d.shiftJISWriter
	}
	return d.defaultWriter
}
