package the_platinum_searcher

import (
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func newEncodeReader(r io.Reader, encoding int) io.Reader {
	switch encoding {
	case EUCJP:
		return transform.NewReader(r, japanese.EUCJP.NewEncoder())
	case SHIFTJIS:
		return transform.NewReader(r, japanese.ShiftJIS.NewEncoder())
	}
	return nil
}
