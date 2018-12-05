package logs

import (
	"archive/zip"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
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
		stat, err := os.Stat(v)
		utils.Check(err)
		if stat.IsDir() {
			utils.Check(i.addDir(v))
		} else {
			utils.Check(i.addZip(v))
		}
	}
}

func idFromFilename(p string) string {
	p = path.Base(p)
	i := strings.Index(p, ".")
	if i == -1 {
		return p
	}
	return p[:i]
}

func addFilesToIndex(index *FileIndex, zipPath string, files []string) error {
	infos := index.CreateZip(zipPath)
	for _, name := range files {
		index.SetFile(infos, idFromFilename(name), name)
	}
	return index.Save()
}

func (i *importer) addDir(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".mjlog.gz") {
			i.index.SetRootFile(idFromFilename(path), path)
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return i.index.Save()
}

func (i *importer) addZip(zipPath string) error {
	return addZipToIndex(i.index, zipPath)
}

func addZipToIndex(index *FileIndex, zipPath string) error {
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
	return addFilesToIndex(index, zipPath, files)
}
