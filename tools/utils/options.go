package utils

import (
	"compress/gzip"
	"io"
	"net/http"
)

type downloader struct {
	tracker    *writeTracker
	compressor func(io.Writer) io.WriteCloser
	client     *http.Client
}

type Option func(*downloader)

func AddTracker(f Tracker) Option {
	return func(x *downloader) {
		x.tracker.add(f)
	}
}

func Client(c *http.Client) Option {
	return func(x *downloader) {
		x.client = c
	}
}

func Compressor(f func(io.Writer) io.WriteCloser) Option {
	return func(x *downloader) {
		x.compressor = f
	}
}

func GZIP() Option {
	return Compressor(func(w io.Writer) io.WriteCloser {
		return gzip.NewWriter(w)
	})
}

func NewDownloader(opts ...Option) *downloader {
	x := &downloader{
		client:  &http.Client{},
		tracker: newWriteTracker(),
	}
	for _, v := range opts {
		v(x)
	}
	return x
}
