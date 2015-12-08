package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func newDecodeReader(r io.Reader, encoding int) io.Reader {
	switch encoding {
	case EUCJP:
		return transform.NewReader(r, japanese.EUCJP.NewDecoder())
	case SHIFTJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	}
	return nil
}
