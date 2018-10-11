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

func (this *AsyncWriter) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for message := range this.writeChan {
		err := this.WriteCallback(ctx, message)
		if err != nil {
			return
		}
	}
}

func (this *AsyncWriter) Close() {
	close(this.writeChan)
}

func (this *AsyncWriter) Start(ctx context.Context) func() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		this.run(ctx)
		wg.Done()
	}()
	return wg.Wait
}

func (this *AsyncWriter) WriteString(message string) error {
	select {
	case this.writeChan <- message:
	default:
		return errors.New("Channel is full")
	}
	return nil
}
