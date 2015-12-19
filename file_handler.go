package the_platinum_searcher

import "os"

func getFileHandler(path string) (*os.File, error) {
	if path == "" {
		return os.Stdin, nil
	} else {
		return os.Open(path)
	}
}
