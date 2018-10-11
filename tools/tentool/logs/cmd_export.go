package logs

import (
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type ConfigRecord struct {
	search *regexp.Regexp
	Search string `json:"search"`
	File   string `json:"file"`
}

type exporter struct {
	interactive bool
	dry         bool
	force       bool
	index       *FileIndex
	configPath  string
	config      []*ConfigRecord
}

func (e *exporter) Run(args []string) {
	index, err := LoadIndex()
	utils.Check(err)
	e.index = index
	utils.Check(e.loadConfig(args[0]))
	mapped := e.mapExportFiles()
	fmt.Printf("Resulting files %v:\n", len(mapped))
	exportCount := 0
	for _, v := range mapped.Keys() {
		count := mapped[v].Count()
		exportCount += count
		fmt.Printf("%v: %v\n", v, count)
	}
	indexCount := 0
	for _, file := range e.index.data {
		if len(file) == 0 {
			continue
		}
		indexCount++
	}
	fmt.Printf("Files to export %v of %v\n", exportCount, indexCount)
	if e.dry {
		return
	}
	for _, v := range mapped.Keys() {
		_, err := os.Stat(v)
		if err == nil && !e.force {
			fmt.Printf("Skipping file %v - already exists\n", v)
			continue
		}
		fmt.Printf("Creating file %v\n", v)
		utils.Check(createZip(v, mapped[v], e.interactive))
	}
}

func getID(p string) string {
	p = path.Base(p)
	i := strings.Index(p, ".")
	if i == -1 {
		return p
	}
	return p[:i]
}

func writeReader(name string, r io.Reader, zw *zip.Writer) error {
	id := getID(name)
	w, err := zw.Create(id + ".mjlog")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

func writeFile(name string, zw *zip.Writer) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	return writeReader(name, gzipReader, zw)
}

func zipIndex(zr *zip.ReadCloser) map[string]*zip.File {
	index := make(map[string]*zip.File, len(zr.File))
	for _, v := range zr.File {
		index[v.Name] = v
	}
	return index
}

func writeFromZip(name string, zr map[string]*zip.File, zw *zip.Writer) error {
	file, err := zr[name].Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return writeReader(name, file, zw)
}

func createZip(name string, desc *exportFile, intractive bool) error {
	iw := utils.NewInteractiveWriter(os.Stdout)
	f, err := utils.CreateFile(name)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(f)
	defer func() {
		zipWriter.Close()
		err = f.CommitOnSuccess(&err)
	}()
	total := desc.Count()
	progress := 0
	onProgress := func() {
		if !intractive {
			return
		}
		progress++
		iw.Printf("%v / %v (%v%%)", progress, total, progress*100/total)
	}
	sort.Strings(desc.files)
	for _, v := range desc.files {
		onProgress()
		err = writeFile(v, zipWriter)
		if err != nil {
			return err
		}
	}
	for _, v := range desc.zipFiles.Keys() {
		var zipReader *zip.ReadCloser
		zipReader, err = zip.OpenReader(v)
		if err != nil {
			return err
		}
		func() {
			defer zipReader.Close()
			files := desc.zipFiles[v]
			sort.Strings(files)
			index := zipIndex(zipReader)
			for _, v := range files {
				onProgress()
				err = writeFromZip(v, index, zipWriter)
				if err != nil {
					return
				}
			}
		}()
		if err != nil {
			return err
		}
	}
	iw.Println()
	return err
}

type exportFile struct {
	files    []string
	zipFiles zipFiles
}

func (f *exportFile) Count() int {
	result := len(f.files)
	for _, v := range f.zipFiles {
		result += len(v)
	}
	return result
}

type zipFiles map[string][]string

type exportFiles map[string]*exportFile

func (f exportFiles) Keys() []string {
	x := make([]string, 0, len(f))
	for k := range f {
		x = append(x, k)
	}
	sort.Strings(x)
	return x
}

func (f exportFiles) get(k string) *exportFile {
	x := f[k]
	if x == nil {
		x = &exportFile{zipFiles: make(zipFiles)}
		f[k] = x
	}
	return x
}

func (f exportFiles) AddFile(k string, args []string) {
	x := f.get(k)
	switch len(args) {
	case 0:
	case 1:
		x.files = append(x.files, args[0])
	case 2:
		x.zipFiles[args[0]] = append(x.zipFiles[args[0]], args[1])
	default:
		panic("Unexpected args len")
	}
}

func (f zipFiles) Keys() []string {
	x := make([]string, 0, len(f))
	for k := range f {
		x = append(x, k)
	}
	sort.Strings(x)
	return x
}

func (e *exporter) mapExportFiles() exportFiles {
	result := exportFiles{}
	for id, file := range e.index.data {
		if len(file) == 0 {
			continue
		}
		for _, c := range e.config {
			if !c.search.MatchString(id) {
				continue
			}
			filename := c.search.ReplaceAllString(id, c.File)
			result.AddFile(filename, file)
			break
		}
	}
	return result
}

func (e *exporter) loadConfig(path string) error {
	var x []*ConfigRecord
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &x)
	if err != nil {
		return err
	}
	for _, v := range x {
		r, err := regexp.Compile(v.Search)
		if err != nil {
			return err
		}
		v.search = r
	}
	e.config = x
	return nil
}
