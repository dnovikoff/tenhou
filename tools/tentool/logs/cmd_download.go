package logs

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

type downloader struct {
	interactive bool
	index       *FileIndex
	parallel    int
	yadiskURL   string
}

func (d *downloader) Run() {
	var err error
	d.index, err = LoadIndex()
	utils.Check(err)
	ids := make([]string, 0, len(d.index.data.Files))
	for _, v := range d.index.data.Files {
		if v.Check() {
			continue
		}
		if v.Failed > 5 {
			continue
		}
		ids = append(ids, v.ID)
	}
	sort.Strings(ids)
	total := len(ids)
	w := utils.NewProgressWriter(os.Stdout, "Downloaded", total).SetDelay(time.Millisecond * 300).SetETA()
	if !d.interactive {
		w = w.Disable()
	}
	fmt.Printf("Logs to download %v of total %v\n", total, d.index.Len())
	ctx := context.TODO()
	requests := make(chan *downloadRequest, d.parallel)
	results := make(chan *downloadRequest, d.parallel*100)

	group := utils.Routines{}
	for i := 0; i < d.parallel; i++ {
		group.Start(func() error {
			d.downloadStream(ctx, requests, results)
			return nil
		})
	}
	go func() {
		group.Wait()
		close(results)
	}()
	schedulerError := make(chan error, 1)
	go d.downloadScheduler(ctx, ids, requests, schedulerError)
	defer d.index.Save()
	w.Start()
	for result := range results {
		w.Inc()
		if result.Error != nil {
			fmt.Printf("\nError downloading %v to %v: %v\n", result.SrcURL, result.DestPath, result.Error)
			d.index.SetError(result.ID)
			continue
		}
		info := d.index.SetRootFile(result.ID, result.DestPath)
		{
			zipped, err := ioutil.ReadFile(result.DestPath)
			utils.Check(err)
			r, err := gzip.NewReader(bytes.NewReader(zipped))
			utils.Check(err)
			data, err := ioutil.ReadAll(r)
			utils.Check(err)
			utils.Check(r.Close())
			names, err := parseNames(data)
			utils.Check(err)
			info.LogNames = names
		}
		if w.Progress()%200 == 0 {
			utils.Check(d.index.Save())
		}
		w.Display()
	}
	w.Done()
	utils.Check(<-schedulerError)
	utils.Check(d.index.Save())
}

type downloadRequest struct {
	SrcURL   string
	DestPath string
	ID       string
	Error    error
}

func (d *downloader) downloadScheduler(ctx context.Context, ids []string, requests chan<- *downloadRequest, resultError chan error) {
	defer close(requests)
	for _, id := range ids {
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
