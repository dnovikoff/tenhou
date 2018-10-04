package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func MustDownload(u string, opts ...Option) string {
	var buf bytes.Buffer
	Check(NewDownloader(opts...).Write(u, &buf))
	return buf.String()
}

func (d *downloader) WriteFile(u, p string) (err error) {
	f, err := CreateFile(p)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
		if err != nil {
			return
		}
		err = f.Commit()
	}()
	err = d.Write(u, f)
	return
}

func (d *downloader) Write(u string, w io.Writer) error {
	resp, err := d.client.Get(u)
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
