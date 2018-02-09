package network

import (
	"context"
	"fmt"
	"io"
)

type xmlDebugger struct {
	impl   XMLConnection
	output io.Writer
}

func NewXMLConnectionDebugger(impl XMLConnection, out io.Writer) XMLConnection {
	return &xmlDebugger{impl, out}
}

func (this *xmlDebugger) Read(ctx context.Context) (str string, err error) {
	str, err = this.impl.Read(ctx)
	if err != nil {
		fmt.Fprintf(this.output, "Get error: %v\n", err)
	} else {
		fmt.Fprintf(this.output, "Get: %v\n", str)
	}
	return
}

func (this *xmlDebugger) Close() error {
	fmt.Fprintf(this.output, "Close\n")
	return this.impl.Close()
}

// Not thread safe
func (this *xmlDebugger) Write(ctx context.Context, str string) error {
	fmt.Fprintf(this.output, "Send: %v\n", str)
	err := this.impl.Write(ctx, str)
	if err != nil {
		fmt.Fprintf(this.output, "Write error: %v\n", err)
	}
	return err
}
