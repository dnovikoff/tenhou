package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tenhou/parser"
)

type nodeReaderTester struct {
	NodeReader
}

func (r *nodeReaderTester) mustError(t *testing.T) {
	_, err := r.Next()
	require.Error(t, err)
}

func (r *nodeReaderTester) mustNext(t *testing.T) *parser.Node {
	n, err := r.Next()
	require.NoError(t, err)
	require.NotNil(t, n)
	return n
}

func newNRTester(messages ...string) *nodeReaderTester {
	return &nodeReaderTester{
		NodeReader{Read: func() (string, error) {
			if len(messages) == 0 {
				return "", fmt.Errorf("No more messages")
			}
			x := messages[0]
			messages = messages[1:]
			return x, nil
		}}}
}

func TestNodeReaderEmpty(t *testing.T) {
	newNRTester().mustError(t)
}

func TestNodeReaderOne(t *testing.T) {
	r := newNRTester("<A/>")
	assert.Equal(t, "A", r.mustNext(t).Name)
	r.mustError(t)
}

func TestNodeReaderMany(t *testing.T) {
	r := newNRTester("<A/><B/><C/>")
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
	assert.Equal(t, "A", r.mustNext(t).Name)
	assert.Equal(t, "B", r.mustNext(t).Name)
	assert.Equal(t, "C", r.mustNext(t).Name)
	assert.Equal(t, "D", r.mustNext(t).Name)
	r.mustError(t)
}

func TestNodeReaderMultiError(t *testing.T) {
	r := newNRTester(
		"<A/>",
		"<B/><C/>",
		"",
		"<<<<>>>>>garbage here",
		"",
		"<D/>")
	assert.Equal(t, "A", r.mustNext(t).Name)
	assert.Equal(t, "B", r.mustNext(t).Name)
	assert.Equal(t, "C", r.mustNext(t).Name)
	r.mustError(t)
}
