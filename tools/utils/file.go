package utils

import (
	"os"
	"path/filepath"
)

func MakeDirForFile(filename string) error {
	dir := filepath.Dir(filename)
	return os.MkdirAll(dir, os.ModePerm)
}
