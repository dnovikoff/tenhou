package logs

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pstats "github.com/dnovikoff/tenhou/genproto/stats"
	"github.com/dnovikoff/tenhou/tools/tentool/stats"
	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

type updater struct {
	interactive bool
	clean       bool
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
	if !u.clean {
		// Ignore error
		u.parsedIndex.Load()
	}
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
	err = stats.ParseStatsForFile(reader, p, func(s *pstats.Record) error {
		u.logsIndex.AddStat(s)
		return nil
	})
	if err != nil {
		fmt.Printf("\nError parsing file %v: %v\n", p, err)
		return err
	}
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
		if v.FileInfo().IsDir() {
			continue
		}
		if !statFileContainsLogs(v.Name) {
			continue
		}
		fileReader, err := v.Open()
		if err != nil {
			return err
		}
		err = func() error {
			defer fileReader.Close()
			reader := fileReader
			if strings.HasSuffix(v.Name, ".gz") {
				gzreader, err := gzip.NewReader(reader)
				if err != nil {
					return err
				}
				reader = gzreader
				defer gzreader.Close()
			}
			err := stats.ParseStatsForFile(reader, v.Name, func(s *pstats.Record) error {
				u.logsIndex.AddStat(s)
				return nil
			})
			if err != nil {
				fmt.Printf("\nError parsing file %v / %v: %v\n", p, v.Name, err)
			}
			return err
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
