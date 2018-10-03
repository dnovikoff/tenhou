package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func MustDownload(u string, opts ...Option) string {
	var buf bytes.Buffer
	Check(NewDownloader(opts...).Write(u, &buf))
	return buf.String()
}

func (d *downloader) WriteFile(u, p string) error {
	err := MakeDirForFile(p)
	if err != nil {
		return err
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	return d.Write(u, f)
}

func (d *downloader) Write(u string, w io.Writer) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code %v on download of %v", resp.StatusCode, u)
	}
	target := w
	if d.compressor != nil {
		c := d.compressor(target)
		defer c.Close()
		target = c
	}
	d.tracker.attach(target, resp.ContentLength)
	return d.tracker.done(io.Copy(d.tracker, resp.Body))
}
