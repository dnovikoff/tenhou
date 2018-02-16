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

func (this *xmlDebugger) Read(ctx context.Context) (str string, err error) {
	str, err = this.impl.Read(ctx)
	if err != nil {
		this.log("Get error: %v", err)
	} else {
		this.log("Get: %v", str)
	}
	return
}

func (this *xmlDebugger) Close() error {
	this.log("Close")
	return this.impl.Close()
}

// Not thread safe
func (this *xmlDebugger) Write(ctx context.Context, str string) error {
	this.log("Send: %v", str)
	err := this.impl.Write(ctx, str)
	if err != nil {
		this.log("Write error: %v", err)
	}
	return err
}
