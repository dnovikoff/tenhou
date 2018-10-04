package logs

import (
	"os"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type FileIndex struct {
	utils.JSONGZFile
	data map[string][]string
}

func NewFileIndex(p string) *FileIndex {
	x := &FileIndex{data: map[string][]string{}}
	x.Path = p
	return x
}

func (i *FileIndex) Load() error {
	return i.JSONGZFile.Load(&i.data)
}

func (i *FileIndex) Save() (err error) {
	return i.JSONGZFile.Save(i.data)
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
