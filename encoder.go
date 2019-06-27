package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func newEncodeReader(r io.Reader, encoding int) io.Reader {
	switch encoding {
	case EUCJP:
		return transform.NewReader(r, japanese.EUCJP.NewEncoder())
	case SHIFTJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	case UTF16LE:
		win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		return transform.NewReader(r, win16le.NewEncoder())
	case UTF16BE:
		win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		return transform.NewReader(r, win16be.NewEncoder())
	}
	return nil
}
