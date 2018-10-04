package utils

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type InteractiveWriter struct {
	w        io.Writer
	buf      *bytes.Buffer
	lastSize int
}

func NewInteractiveWriter(w io.Writer) *InteractiveWriter {
	return &InteractiveWriter{w: w, buf: &bytes.Buffer{}}
}

func (w *InteractiveWriter) Printf(format string, args ...interface{}) (int, error) {
	lastString := w.buf.String()
	w.buf.Reset()
	fmt.Fprintf(w.buf, format, args...)
	if lastString == w.buf.String() {
		return w.buf.Len(), nil
	}
	return fmt.Fprint(w.w, w.fixString(w.buf.String()))
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
