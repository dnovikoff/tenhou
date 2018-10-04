package utils

import (
	"io"
	"os"
	"path/filepath"
)

func MakeDirForFile(filename string) error {
	dir := filepath.Dir(filename)
	return os.MkdirAll(dir, os.ModePerm)
}

type FileWriter struct {
	path    string
	tmpPath string
	file    *os.File
}

var _ io.Writer = &FileWriter{}

func CreateFile(name string) (*FileWriter, error) {
	err := MakeDirForFile(name)
	if err != nil {
		return nil, err
	}
	w := &FileWriter{
		path:    name,
		tmpPath: name + ".tmp",
	}
	w.file, err = os.Create(w.tmpPath)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FileWriter) Write(data []byte) (int, error) {
	return w.file.Write(data)
}

func (w *FileWriter) Commit() error {
	return os.Rename(w.tmpPath, w.path)
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}
