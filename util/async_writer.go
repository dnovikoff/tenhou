package util

import (
	"context"
	"errors"
	"sync"
)

const DefaultChannelSize = 1024

type AsyncWriter struct {
	writeChan     chan string
	WriteCallback func(context.Context, string) error
}

func NewAsyncWriter(channelSize int) *AsyncWriter {
	return &AsyncWriter{writeChan: make(chan string, channelSize)}
}

func (w *AsyncWriter) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for message := range w.writeChan {
		err := w.WriteCallback(ctx, message)
		if err != nil {
			return
		}
	}
}

func (w *AsyncWriter) Close() {
	close(w.writeChan)
}

func (w *AsyncWriter) Start(ctx context.Context) func() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		w.run(ctx)
		wg.Done()
	}()
	return wg.Wait
}

func (w *AsyncWriter) WriteString(message string) error {
	select {
	case w.writeChan <- message:
	default:
		return errors.New("Channel is full")
	}
	return nil
}
