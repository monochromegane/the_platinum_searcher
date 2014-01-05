package file

import (
	"testing"
)

type Assert struct {
	path, fileType string
}

var Asserts = []Assert{
	Assert{"ascii.txt", ASCII},
	Assert{"binary/binary.bin", BINARY},
	Assert{"ja/euc-jp.txt", EUCJP},
	Assert{"ja/shift_jis.txt", SHIFTJIS},
	Assert{"ja/utf8.txt", UTF8},
}

func TestIdentifyType(t *testing.T) {
	for _, f := range Asserts {
		fileType := IdentifyType("../../files/" + f.path)
		if fileType != f.fileType {
			t.Errorf("It should be %s.", f.fileType)
		}
	}
}
