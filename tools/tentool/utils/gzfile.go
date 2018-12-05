package utils

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

type JSONGZFile struct {
	Path   string
	Pretty bool
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
	defer gz.Close()
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
	defer gz.Close()
	enc := json.NewEncoder(gz)
	if f.Pretty {
		enc.SetIndent("", " ")
	}
	err = enc.Encode(data)
	return
}
