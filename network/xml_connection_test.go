package network

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/client"
)

func testPair() (c1, c2 *xmlConnection) {
	server, client := net.Pipe()
	c1 = NewXMLConnection(client)
	c2 = NewXMLConnection(server)
	return
}

func TestConnectionRead(t *testing.T) {
	server, client := testPair()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		require.NoError(t, server.Write(ctx, "Hello World"))
	}()
	data, err := client.Read(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Hello World", data)
}

func TestConnectionReadParts(t *testing.T) {
	server, client := testPair()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		server.impl.Write([]byte("Hello"))
		server.impl.Write([]byte(" World"))
		server.impl.Write([]byte{terminator})
	}()
	data, err := client.Read(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Hello World", data)
}

func TestConnectionReadClosed(t *testing.T) {
	server, client := testPair()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		server.impl.Write([]byte("Hello"))
		server.impl.Write([]byte(" World"))
		server.impl.Close()
	}()
	_, err := client.Read(ctx)
	require.Error(t, err)
}

func TestConnectionSend(t *testing.T) {
	s, c := testPair()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() {
		c := client.NewXMLWriter()
		c.Reach(client.Reach{Step: 1})
		t := client.Take{}
		t.Instance = tile.Pin1.Instance(0)
		c.Take(t)
		s.Write(ctx, c.String())
	}()
	data, err := c.Read(ctx)
	require.NoError(t, err)
	assert.Equal(t, `<REACH who="0" step="1"/> <T36/>`, data)
}
