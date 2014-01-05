package file

import (
	"testing"
)

type FileType struct {
	path, fileType string
}

var Asserts = []FileType{
	FileType{"ascii.txt", ASCII},
	FileType{"binary/binary.bin", BINARY},
	FileType{"ja/euc-jp.txt", EUCJP},
	FileType{"ja/shift_jis.txt", SHIFTJIS},
	FileType{"ja/utf8.txt", UTF8},
}

func TestIdentifyType(t *testing.T) {
	for _, f := range Asserts {
		fileType := IdentifyType("../../files/" + f.path)
		if fileType != f.fileType {
			t.Errorf("It should be %s.", f.fileType)
		}
	}
}
