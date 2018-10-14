package utils

import "io"

type Tracker interface {
	Start(total int64)
	Write(bytes int)
	Done(total int64, err error)
}

func newWriteTracker(f ...Tracker) *writeTracker {
	return &writeTracker{f: f}
}

type writeTracker struct {
	w     io.Writer
	total int64
	f     []Tracker
}

func (w *writeTracker) add(f Tracker) *writeTracker {
	w.f = append(w.f, f)
	return w
}

func (w *writeTracker) attach(x io.Writer, total int64) {
	w.w = x
	if total < 0 {
		total = 0
	}
	for _, v := range w.f {
		v.Start(total)
	}
}

func (w *writeTracker) done(written int64, err error) error {
	for _, v := range w.f {
		v.Done(written, err)
	}
	return err
}

func (w *writeTracker) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	for _, v := range w.f {
		v.Write(n)
	}
	return n, err
}
