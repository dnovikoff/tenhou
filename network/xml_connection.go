package network

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"strings"
)

const terminator = byte(0)

type XMLConnection interface {
	io.Closer
	Read(context.Context) (string, error)
	Write(context.Context, string) error
}

type xmlConnection struct {
	impl   net.Conn
	reader *bufio.Reader
}

var _ XMLConnection = &xmlConnection{}

func NewXMLConnection(impl net.Conn) *xmlConnection {
	return &xmlConnection{
		impl:   impl,
		reader: bufio.NewReader(impl),
	}
}

// Not thread safe
func (c *xmlConnection) Read(ctx context.Context) (str string, err error) {
	type result struct {
		str string
		err error
	}
	ch := make(chan *result, 1)
	go func() {
		var r result
		r.str, r.err = c.reader.ReadString(terminator)
		r.str = strings.TrimRight(r.str, string([]byte{terminator}))
		ch <- &r
	}()
	select {
	case <-ctx.Done():
		err = errors.New("Read timeout exceded")
		c.impl.Close()
	case r := <-ch:
		str, err = r.str, r.err
	}
	return
}

func (c *xmlConnection) Close() error {
	return c.impl.Close()
}

// Not thread safe
func (c *xmlConnection) Write(ctx context.Context, str string) (err error) {
	ch := make(chan error, 1)
	go func() {
		_, writeError := c.impl.Write(append([]byte(str), terminator))
		ch <- writeError
	}()
	select {
	case <-ctx.Done():
		err = errors.New("Write timeout exceded")
		c.impl.Close()
	case err = <-ch:
	}
	return
}
