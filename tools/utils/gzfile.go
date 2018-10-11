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
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	err = json.NewDecoder(gz).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func (f *JSONGZFile) fileName() string {
	return f.Path + ".gz"
}

func (f *JSONGZFile) Save(data interface{}) (err error) {
	file, err := CreateFile(f.fileName())
	if err != nil {
		return
	}
	gz := gzip.NewWriter(file)
	enc := json.NewEncoder(gz)
	err = enc.Encode(data)
	defer func() {
		file.Close()
		if err != nil {
			return
		}
		err = file.Commit()
	}()
	if err != nil {
		return
	}
	err = gz.Close()
	return
}
