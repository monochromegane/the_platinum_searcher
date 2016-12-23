// main package in lib/pt provides search function for C.
// ```
// $ go build -ldflags="-s" -buildmode=c-shared -o libpt.so lib/pt/main.go
// ```
// **THIS IS A BETA FEATURE.**
package main

import (
	"C"
	"os"
	"unsafe"

	pt "github.com/monochromegane/the_platinum_searcher"
)

func goStrings(argc C.int, argv *C.char) []string {
	length := int(argc)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:length:length]
	gostrings := make([]string, length)
	for i, s := range tmpslice {
		gostrings[i] = C.GoString(s)
	}
	return gostrings
}

//export search
func search(s *C.char, optv *C.char, optc C.int) int {
	pt := pt.PlatinumSearcher{Out: os.Stdout, Err: os.Stderr}
	opts := goStrings(optc, optv)
	exitCode := pt.Run(append(opts, C.GoString(s)))
	return exitCode
}

func main() {
}
