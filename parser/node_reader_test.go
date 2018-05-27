package parser

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nodeReaderTester struct {
	NodeReader
	cancel func()
	wait   func()
}

func (r *nodeReaderTester) Stop() {
	r.cancel()
	r.wait()
}

func (r *nodeReaderTester) mustError(t *testing.T) string {
	_, err := r.Next()
	require.Error(t, err)
	return err.Error()
}

func (r *nodeReaderTester) mustNext(t *testing.T) *Node {
	n, err := r.Next()
	require.NoError(t, err)
	require.NotNil(t, n)
	return n
}

func newNRTester(messages ...string) *nodeReaderTester {
	ctx, cancel := context.WithCancel(context.Background())
	x := &nodeReaderTester{NodeReader{}, cancel, nil}
	x.ReadCallback = func(ctx context.Context) (string, error) {
		if len(messages) == 0 {
			return "", fmt.Errorf("No more messages")
		}
		x := messages[0]
		messages = messages[1:]
		return x, nil
	}
	x.wait = x.Start(ctx)
	return x
}

func TestNodeReaderEmpty(t *testing.T) {
	r := newNRTester()
	defer r.Stop()
	r.mustError(t)
}

func TestNodeReaderOne(t *testing.T) {
	r := newNRTester("<A/>")
	defer r.Stop()
	assert.Equal(t, "A", r.mustNext(t).Name)
	r.mustError(t)
}

func TestNodeReaderMany(t *testing.T) {
	r := newNRTester("<A/><B/><C/>")
	defer r.Stop()
	assert.Equal(t, "A", r.mustNext(t).Name)
	assert.Equal(t, "B", r.mustNext(t).Name)
	assert.Equal(t, "C", r.mustNext(t).Name)
	r.mustError(t)
}

func TestNodeReaderMulti(t *testing.T) {
	r := newNRTester(
		"<A/>",
		"<B/><C/>",
		"",
		"",
		"<D/>")
	defer r.Stop()
	assert.Equal(t, "A", r.mustNext(t).Name)
	assert.Equal(t, "B", r.mustNext(t).Name)
	assert.Equal(t, "C", r.mustNext(t).Name)
	assert.Equal(t, "D", r.mustNext(t).Name)
	r.mustError(t)
}

func TestNodeReaderNotStarted(t *testing.T) {
	r := NewNodeReader()
	node, err := r.Next()
	require.Error(t, err)
	require.Equal(t, "NodeReader stopped", err.Error())
	require.Nil(t, node)
}

func TestNodeReaderStopped(t *testing.T) {
	r := newNRTester("<A/>")
	assert.Equal(t, "A", r.mustNext(t).Name)
	r.Stop()
	require.Equal(t, "No more messages", r.mustError(t))
	require.Equal(t, "NodeReader stopped", r.mustError(t))
}

func TestNodeReaderMultiError(t *testing.T) {
	r := newNRTester(
		"<A/>",
		"<B/><C/>",
		"",
		"<<<<>>>>>garbage here",
		"",
		"<D/>")
	defer r.Stop()
	assert.Equal(t, "A", r.mustNext(t).Name)
	assert.Equal(t, "B", r.mustNext(t).Name)
	assert.Equal(t, "C", r.mustNext(t).Name)
	r.mustError(t)
}
