package network

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type xmlDebugger struct {
	impl XMLConnection
	log  func(format string, args ...interface{})
}

func NewMutexLogger(w io.Writer) func(format string, args ...interface{}) {
	var m sync.Mutex
	return func(format string, args ...interface{}) {
		m.Lock()
		defer m.Unlock()
		fmt.Fprintf(w, format+"\n", args...)
	}
}

func NewXMLConnectionDebugger(impl XMLConnection, log func(format string, args ...interface{})) XMLConnection {
	return &xmlDebugger{impl, log}
}

func (d *xmlDebugger) Read(ctx context.Context) (str string, err error) {
	str, err = d.impl.Read(ctx)
	if err != nil {
		d.log("Get error: %v", err)
	} else {
		d.log("Get: %v", str)
	}
	return
}

func (d *xmlDebugger) Close() error {
	d.log("Close")
	return d.impl.Close()
}

// Not thread safe
func (d *xmlDebugger) Write(ctx context.Context, str string) error {
	d.log("Send: %v", str)
	err := d.impl.Write(ctx, str)
	if err != nil {
		d.log("Write error: %v", err)
	}
	return err
}
