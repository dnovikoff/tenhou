package network

import (
	"bufio"
	"context"
	"errors"
	"net"
	"strings"
)

const terminator = byte(0)

type XMLConnection struct {
	impl   net.Conn
	reader *bufio.Reader
}

func NewXMLConnection(impl net.Conn) *XMLConnection {
	return &XMLConnection{
		impl:   impl,
		reader: bufio.NewReader(impl),
	}
}

// Not thread safe
func (this *XMLConnection) Read(ctx context.Context) (str string, err error) {
	type result struct {
		str string
		err error
	}
	ch := make(chan *result, 1)
	go func() {
		var r result
		r.str, r.err = this.reader.ReadString(terminator)
		r.str = strings.TrimRight(r.str, string([]byte{terminator}))
		ch <- &r
	}()
	select {
	case <-ctx.Done():
		err = errors.New("Read timeout exceded")
		this.impl.Close()
	case r := <-ch:
		str, err = r.str, r.err
	}
	return
}

func (this *XMLConnection) Close() {
	this.impl.Close()
}

// Not thread safe
func (this *XMLConnection) Write(ctx context.Context, str string) (err error) {
	ch := make(chan error, 1)
	go func() {
		_, writeError := this.impl.Write(append([]byte(str), terminator))
		ch <- writeError
	}()
	select {
	case <-ctx.Done():
		err = errors.New("Write timeout exceded")
		this.impl.Close()
	case err = <-ch:
	}
	return
}
