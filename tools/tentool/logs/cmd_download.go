package logs

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	"github.com/dnovikoff/tenhou/tools/utils"
)

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
		dst, err := GetFilePath(parsed, id)
		utils.Check(err)
		path := path.Join(Location, dst+".mjlog.gz")
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
	dl := utils.NewDownloader(
		utils.GZIP(),
		utils.Client(d.client),
	)
	err := dl.WriteFile(u, path)
	if err == nil {
		return
	}
	fmt.Printf("Error on downloading %v to %v: %v\n", u, path, err)
	// Dont stop here. Let the other links to be downloaded
}
