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
	"time"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
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
	fmt.Printf("Files to export %v of %v\n", exportCount, index.CountDownloaded())
	if e.dry {
		return
	}
	for _, v := range mapped.Keys() {
		exits, err := utils.FileExists(v)
		utils.Check(err)
		if exits {
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
		if v.FileInfo().IsDir() {
			continue
		}
		index[v.Name] = v
	}
	return index
}

func writeFromZip(f *zip.File, zw *zip.Writer) error {
	file, err := f.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return writeReader(path.Base(f.Name), file, zw)
}

func createZip(name string, desc *exportFile, intractive bool) error {
	f, err := utils.CreateFile(name)
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(f)
	defer func() {
		zipWriter.Close()
		f.CommitOnSuccess(&err)
	}()
	pw := utils.NewProgressWriter(os.Stdout, "", desc.Count()).SetETA().SetDelay(time.Millisecond * 300)
	pw.Start()
	onProgress := func() {
		pw.Inc()
		if !intractive {
			return
		}
		pw.Display()
	}
	for _, v := range desc.Group() {
		if v.IsRoot() {
			for _, f := range v.Files {
				err = writeFile(fileName(f.File), zipWriter)
				if err != nil {
					return err
				}
				onProgress()
			}
		} else {
			checker := make(map[string]bool, len(v.Files))
			for _, f := range v.Files {
				checker[f.File] = true
			}
			var zipReader *zip.ReadCloser
			zipReader, err = zip.OpenReader(fileName(v.File))
			if err != nil {
				return err
			}
			err = func() error {
				defer zipReader.Close()
				for _, f := range zipReader.File {
					if !checker[f.Name] {
						continue
					}
					err = writeFromZip(f, zipWriter)
					if err != nil {
						return err
					}
					onProgress()
				}
				return nil
			}()
			if err != nil {
				return err
			}
		}
	}
	pw.Done()
	return err
}

type exportFile struct {
	infos []*FileInfo
}

func (f *exportFile) Group() []*FileInfos {
	g := make(map[*FileInfos]*FileInfos)
	for _, v := range f.infos {
		x := g[v.parent]
		if x == nil {
			x = &FileInfos{File: v.parent.File, isRoot: v.parent.isRoot}
			g[v.parent] = x
		}
		x.Files = append(x.Files, v)
	}
	out := make([]*FileInfos, 0, len(g))
	for _, v := range g {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].File < out[j].File
	})
	return out
}

func (f *exportFile) Count() int {
	return len(f.infos)
}

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
		x = &exportFile{}
		f[k] = x
	}
	return x
}

func (f *exportFile) add(info *FileInfo) {
	f.infos = append(f.infos, info)
}

func (e *exporter) mapExportFiles() exportFiles {
	out := exportFiles{}
	e.index.data.Visit(func(i *FileInfo) {
		if !i.Check() {
			return
		}
		id := i.ID
		for _, c := range e.config {
			if !c.search.MatchString(id) {
				continue
			}
			filename := c.search.ReplaceAllString(id, c.File)
			out.get(filename).add(i)
			return
		}
	})
	return out
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
