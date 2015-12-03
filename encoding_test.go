package the_platinum_searcher

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

type Assert struct {
	path     string
	fileType int
}

var Asserts = []Assert{
	Assert{"ascii.txt", ASCII},
	Assert{"binary/binary.bin", BINARY},
	Assert{"ja/euc-jp.txt", EUCJP},
	Assert{"ja/shift_jis.txt", SHIFTJIS},
	Assert{"ja/utf8.txt", UTF8},
	Assert{"ja/broken_euc-jp.txt", EUCJP},
	Assert{"ja/broken_shift_jis.txt", SHIFTJIS},
	Assert{"ja/broken_utf8.txt", UTF8},
}

func TestIdentifyType(t *testing.T) {
	for _, f := range Asserts {
		b, _ := ioutil.ReadFile(filepath.Join("files", f.path))
		fileType := detectEncoding(b)
		if fileType != f.fileType {
			t.Errorf("%s should be %d.", f.path, f.fileType)
		}
	}
}
