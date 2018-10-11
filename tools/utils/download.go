package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
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
		err = f.CommitOnSuccess(&err)
	}()
	err = d.Write(u, f)
	return
}

func (d *downloader) Filename(url string) (string, error) {
	resp, err := d.doRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	header := resp.Header.Get("Content-Disposition")
	if header != "" {
		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err != nil {
			return "", err
		}
		filename := params["filename"]
		if filename != "" {
			return filename, nil
		}
	}
	return path.Base(url), nil
}

func (d *downloader) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if d.userAgent != "" {
		req.Header.Set("User-Agent", d.userAgent)
	}
	return req, err
}

func (d *downloader) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := d.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return d.client.Do(req)
}

func (d *downloader) Write(u string, w io.Writer) error {
	resp, err := d.doRequest(http.MethodGet, u, nil)
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
