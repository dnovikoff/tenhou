package logs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type FileIndex struct {
	path string
	data map[string][]string
}

func NewFileIndex(p string) *FileIndex {
	return &FileIndex{path: p, data: map[string][]string{}}
}

func (i *FileIndex) MakeDir() error {
	return utils.MakeDirForFile(i.fileName())
}

func (i *FileIndex) Load() error {
	f, err := os.Open(i.fileName())
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	err = json.NewDecoder(gz).Decode(&i.data)
	if err != nil {
		return err
	}
	return nil
}

func (i *FileIndex) fileName() string {
	return i.path + ".gz"
}

func (i *FileIndex) Save() error {
	buf := &bytes.Buffer{}
	gz := gzip.NewWriter(buf)
	enc := json.NewEncoder(gz)
	err := enc.Encode(i.data)
	if err != nil {
		return err
	}
	err = gz.Close()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(i.fileName(), buf.Bytes(), 0644)
}

func (i *FileIndex) Len() int {
	return len(i.data)
}

func (i *FileIndex) Add(ids []string) {
	if len(ids) == 0 {
		return
	}
	for _, v := range ids {
		_, found := i.data[v]
		if !found {
			i.data[v] = nil
		}
	}
	return
}

func (i *FileIndex) Set(id string, path []string) {
	i.data[id] = path
}

func (i *FileIndex) Get(id string) []string {
	return i.data[id]
}

func (i *FileIndex) Validate() error {
	var total error
	checked := make(map[string]bool, len(i.data))
	for k, path := range i.data {
		if len(path) == 0 {
			continue
		}
		// TODO: check subs
		v := path[0]
		found, ok := checked[v]
		if !ok {
			_, err := os.Stat(v)
			if err != nil {
				total = multierror.Append(total, err)
				checked[v] = false
				found = false
			} else {
				checked[v] = true
				found = true
			}
		}
		if !found {
			delete(i.data, k)
		}
	}
	return total
}
