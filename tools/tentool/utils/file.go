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

func (w *FileWriter) CommitOnSuccess(err *error) {
	w.Close()
	if *err != nil {
		return
	}
	*err = w.Commit()
}

func (w *FileWriter) Commit() error {
	return os.Rename(w.tmpPath, w.path)
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
