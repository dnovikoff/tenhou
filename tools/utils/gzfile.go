package utils

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

type JSONGZFile struct {
	Path string
}

func (f *JSONGZFile) Load(out interface{}) error {
	file, err := os.Open(f.fileName())
	if err != nil {
		return err
	}
	gz, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return err
	}
	defer func() {
		gz.Close()
		file.Close()
	}()
	return json.NewDecoder(gz).Decode(out)
}

func (f *JSONGZFile) fileName() string {
	return f.Path + ".gz"
}

func (f *JSONGZFile) Save(data interface{}) (err error) {
	file, err := CreateFile(f.fileName())
	if err != nil {
		return
	}
	defer file.CommitOnSuccess(&err)
	gz := gzip.NewWriter(file)
	enc := json.NewEncoder(gz)
	err = enc.Encode(data)
	if err != nil {
		return
	}
	err = gz.Close()
	return
}
