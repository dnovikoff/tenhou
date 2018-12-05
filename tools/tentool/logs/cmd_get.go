package logs

import (
	"archive/zip"
	"compress/gzip"
	"io"
	"log"
	"os"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

func cmdGet(id string) {
	parsedIds := ParseIDs(id)
	if len(parsedIds) == 1 {
		id = parsedIds[0]
	}
	index, err := LoadIndex()
	utils.Check(err)
	info := index.Get(id)
	if info.IsInsideZip() {
		zf, err := zip.OpenReader(fileName(info.parent.File))
		utils.Check(err)
		defer zf.Close()
		for _, v := range zf.File {
			if v.Name != info.File {
				continue
			}
			f, err := v.Open()
			utils.Check(err)
			defer f.Close()
			io.Copy(os.Stdout, f)
			return
		}
		log.Fatalf("File '%v' not found in archive '%v'", info.File, info.parent.File)
	} else if info.File != "" {
		f, err := os.Open(fileName(info.File))
		utils.Check(err)
		gz, err := gzip.NewReader(f)
		utils.Check(err)
		io.Copy(os.Stdout, gz)
		utils.Check(gz.Close())
		utils.Check(f.Close())
	} else {
		log.Fatalf("Log with id '%v' not found in database", id)
	}
}
