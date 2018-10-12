package logs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/dnovikoff/tenhou/tools/utils"
)

type downloader struct {
	interactive bool
	index       *FileIndex
	parallel    int
}

func (d *downloader) Run() {
	var err error
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
	ctx := context.TODO()
	requests := make(chan *downloadRequest, d.parallel)
	results := make(chan *downloadRequest, d.parallel)
	var wg sync.WaitGroup
	for i := 0; i < d.parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.downloadStream(ctx, requests, results)
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	schedulerError := make(chan error, 1)
	go d.downloadScheduler(ctx, links, requests, schedulerError)
	defer d.index.Save()
	for result := range results {
		progress++
		if result.Error != nil {
			fmt.Printf("\nError downloading %v to %v: %v\n", result.SrcURL, result.DestPath, result.Error)
			continue
		}
		d.index.Set(result.ID, []string{result.DestPath})
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
	utils.Check(<-schedulerError)
}

type downloadRequest struct {
	SrcURL   string
	DestPath string
	ID       string
	Error    error
}

func (d *downloader) downloadScheduler(ctx context.Context, links []string, requests chan<- *downloadRequest, resultError chan error) {
	defer close(requests)
	for _, id := range links {
		parsed := ParseID(id)
		downloadLink := GetDownloadLink(parsed)
		dst, err := GetFilePath(parsed, id)
		if err != nil {
			resultError <- err
			return
		}
		path := path.Join(Location, dst+".mjlog.gz")
		select {
		case requests <- &downloadRequest{
			SrcURL:   downloadLink,
			DestPath: path,
			ID:       id,
		}:
		case <-ctx.Done():
			return
		}
	}
	resultError <- nil
}

func (d *downloader) downloadStream(ctx context.Context, input <-chan *downloadRequest, output chan<- *downloadRequest) {
	dl := utils.NewDownloader(
		utils.GZIP(),
		utils.Client(&http.Client{}),
	)
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-input:
			if req == nil {
				return
			}
			req.Error = dl.WriteFile(ctx, req.SrcURL, req.DestPath)
			output <- req
		}
	}
}
