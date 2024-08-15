package lib

import "path/filepath"

func PathsSplit(filename string) (string, string, string) {
	dir, file := filepath.Split(filename)
	ext := filepath.Ext(file)
	rest := len(file) - len(ext)

	if rest == 0 {
		return dir, file, ""
	}

	return dir, file[:rest], ext
}
