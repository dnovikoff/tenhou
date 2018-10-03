package stats

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type FileIndex struct {
	path string
	data map[string]string
}

func NewFileIndex(p string) *FileIndex {
	return &FileIndex{path: p, data: map[string]string{}}
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

func (i *FileIndex) JustAdd(url, path string) {
	i.data[url] = path
}

func (i *FileIndex) Add(url, path string) error {
	i.JustAdd(url, path)
	return i.Save()
}

func (i *FileIndex) Check(url string) bool {
	_, found := i.data[url]
	return found
}

func (i *FileIndex) Len() int {
	return len(i.data)
}

func (i *FileIndex) Files() []string {
	out := make([]string, 0, len(i.data))
	for _, v := range i.data {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func (i *FileIndex) Validate() error {
	var total error
	checked := make(map[string]bool, len(i.data))
	for k, v := range i.data {
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
