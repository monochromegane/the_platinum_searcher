package pt

import (
	"fmt"
	"os"
	"strings"
)

const (
	BINARY   = "Binary"
	ASCII    = "ASCII"
	UTF8     = "UTF-8"
	EUCJP    = "EUC-JP"
	SHIFTJIS = "Shift_JIS"
)

func IdentifyFileType(path string) (fileType string) {

	var (
		suspiciousBytes = 0
		likelyUtf8      = 0
		likelyEucjp     = 0
		likelyShiftjis  = 0
	)

	file, _ := os.Open(path)
	defer file.Close()

	stat, _ := file.Stat()
	bs := make([]byte, stat.Size())
	file.Read(bs)

	var length = len(bs)
	var total = length
	if total > 512 {
		total = 512
	}

	if length == 0 {
		return ASCII
	}

	if length >= 3 && bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		// UTF-8 BOM. This isn't binary.
		return UTF8
	}

	if length >= 5 && bs[0] == 0x25 && bs[1] == 0x50 && bs[2] == 0x44 && bs[3] == 0x46 && bs[4] == 0x2D {
		/*  %PDF-. This is binary. */
		return BINARY
	}

	for i := 0; i < total; i++ {
		if bs[i] == 0x00 {
			/* NULL char. It's binary */
			return BINARY
		} else if (bs[i] < 7 || bs[i] > 14) && (bs[i] < 32 || bs[i] > 127) {
			/* UTF-8 detection */
			if bs[i] > 193 && bs[i] < 224 && i+1 < total {
				i++
				if bs[i] > 127 && bs[i] < 192 {
					likelyUtf8++
					continue
				}

			} else if bs[i] > 223 && bs[i] < 240 && i+2 < total {
				i++
				if bs[i] > 127 && bs[i] < 192 && bs[i+1] > 127 && bs[i+1] < 192 {
					i++
					likelyUtf8++
					continue
				}
			}

			/* EUC-JP detection */
			if bs[i] == 142 {
				i++
				if bs[i] > 160 && bs[i] < 224 {
					likelyEucjp++
					continue
				}
			} else if bs[i] > 160 && bs[i] < 255 {
				i++
				if bs[i] > 160 && bs[i] < 255 {
					likelyEucjp++
					continue
				}
			}

			/* Shift-JIS detection */
			if bs[i] > 160 && bs[i] < 224 {
				likelyShiftjis++
				continue
			} else if (bs[i] > 128 && bs[i] < 160) || (bs[i] > 223 && bs[i] < 240) {
				i++
				if (bs[i] > 63 && bs[i] < 127) || (bs[i] > 127 && bs[i] < 253) {
					likelyShiftjis++
					continue
				}
			}

			suspiciousBytes++
			if i >= 32 && (suspiciousBytes*100)/total > 10 {
				return BINARY
			}

		}
	}

	if (suspiciousBytes*100)/total > 10 {
		return BINARY
	}

	// fmt.Printf("Detected points[utf8/eucjp/shiftjis] is %d/%d/%d.\n", likelyUtf8, likelyEucjp, likelyShiftjis)

	if likelyUtf8 == 0 && likelyEucjp == 0 && likelyShiftjis == 0 {
		return ASCII
	} else if likelyUtf8 >= likelyEucjp && likelyUtf8 >= likelyShiftjis {
		return UTF8
	} else if likelyEucjp >= likelyUtf8 && likelyEucjp >= likelyShiftjis {
		return EUCJP
	} else if likelyShiftjis >= likelyUtf8 && likelyShiftjis >= likelyEucjp {
		return SHIFTJIS
	}

	return ASCII

}

const (
	ColorReset      = "\x1b[0m\x1b[K"
	ColorLineNumber = "\x1b[1;33m"  /* yellow with black background */
	ColorPath       = "\x1b[1;32m"  /* bold green */
	ColorMatch      = "\x1b[30;43m" /* black with yellow background */
)

func PrintPath(path string) {
	fmt.Printf("%s%s%s:", ColorPath, path, ColorReset)
}
func PrintLineNumber(lineNum int) {
	fmt.Printf("%s%d%s:", ColorLineNumber, lineNum, ColorReset)
}
func PrintMatch(pattern, match string) {
	fmt.Printf("%s", strings.Replace(match, pattern, ColorMatch + pattern + ColorReset, -1))
}
