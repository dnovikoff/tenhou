package logs

import (
	"archive/zip"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type importer struct {
	index       *FileIndex
	interactive bool
}

func (i *importer) Run(args []string) {
	var err error
	i.index, err = LoadIndex()
	utils.Check(err)
	for _, v := range args {
		utils.Check(i.addZip(v))
	}
}

func (i *importer) addToIndex(zipPath string, files []string) error {
	for _, name := range files {
		ext := path.Ext(name)
		id := strings.TrimSuffix(path.Base(name), ext)
		i.index.Set(id, []string{zipPath, name})
	}
	return i.index.Save()
}

func (i *importer) addZip(zipPath string) error {
	indexPath, err := makeZipPath(zipPath)
	if err != nil {
		return err
	}
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	files := make([]string, 0, len(reader.File))
	for _, v := range reader.File {
		if v.FileInfo().IsDir() {
			continue
		}
		files = append(files, v.Name)
	}
	fmt.Printf("Adding %v logs from %v\n", len(files), zipPath)
	return i.addToIndex(indexPath, files)
}

func makeZipPath(p string) (string, error) {
	p = filepath.Clean(p)
	if filepath.IsAbs(p) {
		return p, nil
	}
	return filepath.Rel(Location, p)
}
