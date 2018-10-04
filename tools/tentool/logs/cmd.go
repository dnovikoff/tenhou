package logs

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dnovikoff/tenhou/tools/tentool/stats"
	"github.com/dnovikoff/tenhou/tools/utils"
	"github.com/spf13/cobra"
)

const (
	Location = "./tenhou/logs/"
)

func NewIndex() *FileIndex {
	return NewFileIndex(filepath.Join(Location, "index.json"))
}

func NewParsedIndex() *stats.FileIndex {
	return stats.NewFileIndex(filepath.Join(Location, "parsed.json"))
}

func LoadIndex() (*FileIndex, error) {
	i := NewIndex()
	err := i.Load()
	if err != nil {
		return nil, err
	}
	return i, err
}

func CMD() *cobra.Command {
	rootCMD := &cobra.Command{
		Use:   "logs",
		Short: "Work with tenhou logs",
	}
	initCMD := &cobra.Command{
		Use:   "init",
		Short: "Initialize database in current folder",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			index := NewIndex()
			err := index.Load()
			if os.IsNotExist(err) {
				utils.Check(index.Save())
				fmt.Println("Index initialized for ", index.Path)
				return
			}
			if err == nil {
				fmt.Println("Index already initialized. Delete " + index.Path + " to reinit.")
				os.Exit(1)
			}
		},
	}
	var repairFlag bool
	validateCMD := &cobra.Command{
		Use:   "validate",
		Short: "Validate index file, and existance of mentioned files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			index, err := LoadIndex()
			utils.Check(err)
			err = index.Validate()
			if err == nil {
				fmt.Println("Index is valid")
				return
			}
			if err != nil {
				if repairFlag {
					utils.Check(index.Save())
					fmt.Println("Index repaired")
				} else {
					utils.Check(err)
				}
			}
		},
	}
	validateCMD.Flags().BoolVar(&repairFlag, "repair", false, "Repair index if broken. Delete all records of not found files.")

	var interactiveFlag bool
	updateCMD := &cobra.Command{
		Use:   "update",
		Short: "Update links index from downloaded stat files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			(&updater{interactive: interactiveFlag}).Run()
		},
	}
	updateCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive downloader progress.")

	downloadCMD := &cobra.Command{
		Use:   "download",
		Short: "Download log files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			(&downloader{interactive: interactiveFlag}).Run()
		},
	}
	downloadCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive downloader progress.")

	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(validateCMD)
	rootCMD.AddCommand(updateCMD)
	rootCMD.AddCommand(downloadCMD)
	return rootCMD
}

type downloader struct {
	interactive bool
	index       *FileIndex
	client      *http.Client
}

func (d *downloader) Run() {
	var err error
	d.client = &http.Client{}
	d.index, err = LoadIndex()
	utils.Check(err)
	links := make([]string, 0, len(d.index.data))
	for k, v := range d.index.data {
		if v != nil {
			continue
		}
		links = append(links, k)
	}
	sort.Strings(links)
	total := len(links)
	progress := 0
	w := utils.NewInteractiveWriter(os.Stdout)
	w.Printf("Logs to download %v of total %v", total, len(d.index.data))
	w.Println()
	startTime := time.Now()
	defer d.index.Save()
	for _, id := range links {
		parsed := ParseID(id)
		downloadLink := GetDownloadLink(parsed)
		dst, err := GetFilePath(parsed)
		utils.Check(err)
		path := Location + "/" + dst + ".xml.gz"
		d.download(downloadLink, path)
		d.index.Set(id, []string{path})
		progress++
		if progress%50 == 0 {
			utils.Check(d.index.Save())
		}
		if d.interactive {
			currentTime := time.Now()
			elapsed := currentTime.Sub(startTime)
			itemsLeft := total - progress
			var speed float64
			nanos := elapsed.Nanoseconds()
			if nanos != 0 {
				speed = float64(elapsed.Nanoseconds()) / float64(progress)
			}
			left := time.Nanosecond * time.Duration(speed*float64(itemsLeft))
			left = left.Truncate(time.Second)
			w.Printf("Downloaded %v/%v (%v%%) Time left: %v", progress, total, progress*100/total, left)
		}
	}
	fmt.Println()
}

func (d *downloader) download(u, path string) {
	utils.Check(
		utils.NewDownloader(
			utils.GZIP(),
			utils.Client(d.client),
		).WriteFile(u, path),
	)
}

type updater struct {
	interactive bool
	logsIndex   *FileIndex
	parsedIndex *stats.FileIndex
}

func (u *updater) Run() {
	var err error
	u.logsIndex, err = LoadIndex()
	utils.Check(err)
	statsIndex, err := stats.LoadIndex()
	utils.Check(err)
	u.parsedIndex = NewParsedIndex()
	// Ignore error
	u.parsedIndex.Load()
	size := statsIndex.Len()
	progress := 0
	oldLinksSize := u.logsIndex.Len()
	w := utils.NewInteractiveWriter(os.Stdout)
	w.Printf("Files in index %v", size)
	w.Println()
	for _, v := range statsIndex.Files() {
		progress++
		if u.interactive {
			w.Printf("Parsed %v of %v files (%v%%)", progress, size, progress*100/size)
		}
		utils.Check(u.parseFile(v))
	}
	w.Println()
	utils.Check(u.save())
	newLinksSize := u.logsIndex.Len()
	w.Printf("Parsed %v new ids. Total ids in database %v\n", newLinksSize-oldLinksSize, newLinksSize)
	w.Println()
}

func (u *updater) parseFile(p string) error {
	if u.parsedIndex.Check(p) {
		return nil
	}
	if strings.HasSuffix(p, ".gz") {
		return u.parseGZ(p)
	} else if strings.HasSuffix(p, ".zip") {
		return u.parseZip(p)
	}
	return nil
}

func (u *updater) parseGZ(p string) error {
	if !statFileContainsLogs(p) {
		return nil
	}
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	reader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	u.logsIndex.Add(ParseIDs(string(bytes)))
	u.parsedIndex.JustAdd(p, "")
	return nil
}

func (u *updater) parseZip(p string) error {
	reader, err := zip.OpenReader(p)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, v := range reader.File {
		if !statFileContainsLogs(v.Name) {
			continue
		}
		fileReader, err := v.Open()
		if err != nil {
			return err
		}
		err = func() error {
			defer fileReader.Close()
			bytes, err := ioutil.ReadAll(fileReader)
			if err != nil {
				return err
			}
			u.logsIndex.Add(ParseIDs(string(bytes)))
			return nil
		}()
		if err != nil {
			return err
		}
	}
	u.parsedIndex.JustAdd(p, "")
	return u.save()
}

func (u *updater) save() error {
	err := u.logsIndex.Save()
	if err != nil {
		return err
	}
	return u.parsedIndex.Save()
}

func statFileContainsLogs(p string) bool {
	filename := filepath.Base(p)
	return strings.HasPrefix(filename, "scc") ||
		strings.HasPrefix(filename, "sce") ||
		strings.HasPrefix(filename, "scf")
}
