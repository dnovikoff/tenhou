package util

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

type awTester struct {
	AsyncWriter
	Messages []string
}

func newAwTester() *awTester {
	tst := &awTester{}
	tst.ChannelSize = 10
	tst.WriteCallback = func(_ context.Context, message string) (err error) {
		if message == "error" {
			return errors.New("You wanted error - here you are")
		}
		tst.Messages = append(tst.Messages, message)
		return
	}
	return tst
}

func (this *awTester) start() func() {
	return this.Start(context.Background())
}

func TestAsyncWriter(t *testing.T) {
	tst := newAwTester()
	w := tst.start()
	tst.Close()
	w()
	assert.Empty(t, tst.Messages)
}

func TestAsyncWriterOne(t *testing.T) {
	tst := newAwTester()
	w := tst.start()
	require.NoError(t, tst.WriteString("Hello!"))
	tst.Close()
	w()
	assert.Equal(t, []string{
		"Hello!",
	}, tst.Messages)
}

func TestAsyncWriterTwo(t *testing.T) {
	tst := newAwTester()
	w := tst.start()
	require.NoError(t, tst.WriteString("Hello!"))
	require.NoError(t, tst.WriteString("World!"))
	tst.Close()
	w()
	assert.Equal(t, []string{
		"Hello!",
		"World!",
	}, tst.Messages)
}

func TestAsyncWriterError(t *testing.T) {
	tst := newAwTester()
	w := tst.start()
	require.NoError(t, tst.WriteString("One"))
	require.NoError(t, tst.WriteString("error"))
	require.NoError(t, tst.WriteString("Two"))
	tst.Close()
	w()
	assert.Equal(t, []string{
		"One",
	}, tst.Messages)
}

func TestAsyncWriterQueueIfFulll(t *testing.T) {
	tst := newAwTester()
	tst.ChannelSize = 1
	w := tst.start()
	require.NoError(t, tst.WriteString("error"))
	require.Error(t, tst.WriteString("One"))
	tst.Close()
	w()
	assert.Nil(t, tst.Messages)
}
