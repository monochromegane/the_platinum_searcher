package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func newDecodeReader(r io.Reader, encoding int) io.Reader {
	switch encoding {
	case EUCJP:
		return transform.NewReader(r, japanese.EUCJP.NewDecoder())
	case SHIFTJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	case UTF16LE:
		win16le := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		return transform.NewReader(r, win16le.NewDecoder())
	case UTF16BE:
		win16be := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		return transform.NewReader(r, win16be.NewDecoder())
	}
	return nil
}
