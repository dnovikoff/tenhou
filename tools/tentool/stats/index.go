package stats

import (
	"os"
	"sort"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type FileIndex struct {
	utils.JSONGZFile
	data map[string]string
}

func NewFileIndex(p string) *FileIndex {
	x := &FileIndex{data: map[string]string{}}
	x.Path = p
	return x
}

func (i *FileIndex) Load() error {
	return i.JSONGZFile.Load(&i.data)
}

func (i *FileIndex) Save() (err error) {
	return i.JSONGZFile.Save(i.data)
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
