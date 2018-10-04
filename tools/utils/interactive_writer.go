package utils

import (
	"fmt"
	"io"
	"strings"
)

type InteractiveWriter struct {
	w        io.Writer
	lastSize int
}

func NewInteractiveWriter(w io.Writer) *InteractiveWriter {
	return &InteractiveWriter{w: w}
}

func (w *InteractiveWriter) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(w.w, w.fixString(format), args...)
}

func (w *InteractiveWriter) Println(args ...interface{}) (int, error) {
	w.lastSize = 0
	return fmt.Fprintln(w.w, args...)
}

func (w *InteractiveWriter) fixString(str string) string {
	size := len(str)
	if w.lastSize != 0 {
		str = "\r" + str
	}
	diff := w.lastSize - size
	if diff > 0 {
		str += strings.Repeat(" ", diff)
	}
	w.lastSize = size
	return str
}
