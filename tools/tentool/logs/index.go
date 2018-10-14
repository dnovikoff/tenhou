package logs

import (
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"go.uber.org/multierr"

	pstats "github.com/dnovikoff/tenhou/genproto/stats"
	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

type FileInfo struct {
	ID        string   `json:"id"`
	File      string   `json:"file,omitempty"`
	Failed    int      `json:"failed,omitempty"`
	LogNames  []string `json:"log-names,omitempty"`
	StatNames []string `json:"stat-names,omitempty"`

	// Skip for encode, recover at decode
	parent *FileInfos
	// Special flag for cleanup
	Remove bool `json:"remove,omitempty"`
}

type FileInfos struct {
	File  string       `json:"file,omitempty"`
	Files []*FileInfo  `json:"files,omitempty"`
	Zips  []*FileInfos `json:"zips,omitempty"`

	isRoot bool
}

func (i *FileInfos) Count() int {
	cnt := 0
	for _, v := range i.Zips {
		cnt += v.Count()
	}
	return cnt + len(i.Files)
}

func (i *FileInfo) IsInsideZip() bool {
	return !i.parent.IsRoot()
}

func (i *FileInfos) Visit(f func(info *FileInfo)) {
	for _, v := range i.Files {
		f(v)
	}
	for _, v := range i.Zips {
		v.Visit(f)
	}
}

func (i *FileInfos) IsRoot() bool {
	return i == nil || i.isRoot
}

func (i *FileInfo) Check() bool {
	if i == nil {
		return false
	}
	return !i.Remove && i.File != ""
}

type FileInfoMap map[string]*FileInfo
type ZipIndex map[string]*FileInfos

func (f *FileInfos) optimize() {
	i := 0
	for _, v := range f.Files {
		if !v.Remove {
			f.Files[i] = v
			i++
		}
	}
	f.Files = f.Files[:i]
	i = 0
	for _, v := range f.Zips {
		v.optimize()
		if !v.Empty() {
			f.Zips[i] = v
			i++
		}
	}
	f.Zips = f.Zips[:i]
}

func (f *FileInfos) Empty() bool {
	if f == nil {
		return true
	}
	return len(f.Files) == 0 && len(f.Zips) == 0
}

func (f *FileInfos) zipIndex() ZipIndex {
	x := make(ZipIndex, len(f.Zips))
	for _, v := range f.Zips {
		x[v.File] = v
	}
	return x
}

func (f *FileInfos) index(mp FileInfoMap) {
	for _, v := range f.Files {
		if old := mp[v.ID]; old != nil {
			old.Remove = true
		}
		mp[v.ID] = v
		v.parent = f
	}
	for _, v := range f.Zips {
		v.index(mp)
	}
}

func (f *FileInfos) Index() FileInfoMap {
	mp := make(FileInfoMap, len(f.Files))
	f.index(mp)
	f.isRoot = true
	return mp
}

type FileIndex struct {
	utils.JSONGZFile
	data     *FileInfos
	indexed  FileInfoMap
	zipIndex ZipIndex
}

func NewFileIndex(p string) *FileIndex {
	x := &FileIndex{
		data:    &FileInfos{},
		indexed: FileInfoMap{},
	}
	x.Path = p
	x.Pretty = true
	return x
}

func (i *FileIndex) Load() error {
	var data FileInfos
	err := i.JSONGZFile.Load(&data)
	if err != nil {
		return err
	}
	i.data = &data
	i.reindex()
	return nil
}

func (i *FileIndex) reindex() {
	i.indexed = i.data.Index()
	i.zipIndex = i.data.zipIndex()
}

func (i *FileIndex) Save() error {
	i.data.optimize()
	err := i.JSONGZFile.Save(i.data)
	go runtime.GC()
	return err
}

func (i *FileIndex) Len() int {
	return len(i.indexed)
}

func (i *FileIndex) CountDownloaded() int {
	x := 0
	for _, v := range i.indexed {
		if v.Check() {
			x++
		}
	}
	return x
}

func (i *FileIndex) CreateZip(file string) *FileInfos {
	file = importLocation(file)
	zip := i.zipIndex[file]
	if zip == nil {
		zip = &FileInfos{File: file}
		i.zipIndex[file] = zip
		i.data.Zips = append(i.data.Zips, zip)
	} else {
		zip.Files = nil
		zip.Zips = nil
		i.reindex()
	}
	return zip
}

func (i *FileIndex) AddStat(v *pstats.Record) bool {
	isNew := false
	info := i.Get(v.Id)
	if info == nil {
		isNew = true
		info = i.SetRootFile(v.Id, "")
	}
	names := make([]string, 0, 4)
	for _, p := range v.Players {
		names = append(names, p.Name)
	}
	info.StatNames = names
	return isNew
}

func (i *FileIndex) AddIDs(ids []string) int {
	count := 0
	for _, v := range ids {
		if i.Get(v) == nil {
			count++
			i.SetRootFile(v, "")
		}
	}
	return count
}

func importLocation(path string) string {
	return strings.TrimPrefix(path, Location)
}

func (i *FileIndex) SetFile(parent *FileInfos, id string, file string) *FileInfo {
	info := i.Get(id)
	file = importLocation(file)
	if info == nil {
		info = &FileInfo{ID: id, File: file}
	} else {
		cp := *info
		info.Remove = true
		cp.File = file
		cp.parent = parent
		cp.LogNames = nil
		cp.parent = parent
		info = &cp
	}
	parent.Files = append(parent.Files, info)
	i.indexed[id] = info
	return info
}

func (i *FileIndex) SetError(id string) *FileInfo {
	fi := i.Get(id)
	if fi != nil {
		fi.Failed++
	}
	return fi
}

func (i *FileIndex) SetRootFile(id string, file string) *FileInfo {
	return i.SetFile(i.data, id, file)
}

func (i *FileIndex) Get(id string) *FileInfo {
	return i.indexed[id]
}

func (i *FileIndex) Check(id string) bool {
	return i.indexed[id].Check()
}

type Opener interface {
	Open() (io.ReadCloser, error)
}

type gzFileOpener struct {
	name string
	f    io.Closer
	gz   io.ReadCloser
}

func (o *gzFileOpener) Close() error {
	return multierr.Append(
		o.f.Close(),
		o.gz.Close(),
	)
}

func (o *gzFileOpener) Read(p []byte) (int, error) {
	return o.gz.Read(p)
}

func (o *gzFileOpener) Open() (io.ReadCloser, error) {
	f, err := os.Open(o.name)
	if err != nil {
		return nil, err
	}
	gz, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	o.f, o.gz = f, gz
	return o, nil
}

func fileName(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return path.Join(Location, p)
}

func (i *FileIndex) Visit(f func(info *FileInfo, opener Opener) error) error {
	for _, v := range i.data.Files {
		if !v.Check() {
			continue
		}
		err := f(v, &gzFileOpener{name: fileName(v.File)})
		if err != nil {
			return err
		}
	}
	for _, z := range i.data.Zips {
		zr, err := zip.OpenReader(fileName(z.File))
		if err != nil {
			return err
		}
		filesToCheck := make(map[string]*FileInfo, len(z.Files))
		for _, v := range z.Files {
			filesToCheck[v.File] = v
		}
		for _, v := range zr.File {
			info := filesToCheck[v.Name]
			if !info.Check() {
				continue
			}
			err := f(info, v)
			if err != nil {
				return err
			}
		}
		err = zr.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func sortedNames(x []string) []string {
	var out []string
	out = append(out, x...)
	sort.Strings(out)
	return out
}

func equalNames(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	x, y = sortedNames(x), sortedNames(y)
	for k, v := range x {
		if v != y[k] {
			return false
		}
	}
	return true
}

func checkNames(x, y []string) bool {
	// Do not delete if info is not full
	if len(x) == 0 || len(y) == 0 {
		return true
	}
	return equalNames(x, y)
}

func (i *FileIndex) ValidateNames() []*FileInfo {
	var out []*FileInfo
	for _, v := range i.indexed {
		v.StatNames = fixNames(v.StatNames)
		v.LogNames = fixNames(v.LogNames)
		if !checkNames(v.StatNames, v.LogNames) {
			info := i.SetRootFile(v.ID, "")
			info.StatNames = v.StatNames
			out = append(out, v)
		}
	}
	return out
}

func (i *FileIndex) Validate() ([]string, error) {
	var removed []string
	filesToCheck := make([]string, 0, len(i.data.Files)+len(i.data.Zips))
	for _, v := range i.data.Files {
		filesToCheck = append(filesToCheck, v.File)
	}
	for _, v := range i.data.Zips {
		filesToCheck = append(filesToCheck, v.File)
	}
	for _, v := range filesToCheck {
		exists, err := utils.FileExists(fileName(v))
		if err != nil {
			return nil, err
		}
		if !exists {
			removed = append(removed, v)
		}
	}
	return removed, nil
}
