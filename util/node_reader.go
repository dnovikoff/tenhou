package util

import (
	"context"
	"errors"
	"sync"

	"github.com/dnovikoff/tenhou/parser"
)

type readResult struct {
	node parser.Node
	err  error
}

type NodeReader struct {
	resultCh chan readResult

	wg sync.WaitGroup
}

func NewNodeReader() *NodeReader {
	ch := make(chan readResult)
	// Create with closed channel for correct errors
	close(ch)
	return &NodeReader{resultCh: ch}

}

func (this *NodeReader) run(ctx context.Context, readF func(context.Context) (string, error)) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		this.wg.Done()
		close(this.resultCh)
	}()
	for {
		message, err := readF(ctx)
		if err != nil {
			this.resultCh <- readResult{err: err}
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
		nodes, err := ParseXML(message)
		if err != nil {
			this.resultCh <- readResult{err: err}
			return
		}
		for _, v := range nodes {
			this.resultCh <- readResult{v, nil}
		}
	}
}

func (this *NodeReader) Wait() {
	this.wg.Wait()
}

func (this *NodeReader) Start(ctx context.Context, read func(context.Context) (string, error)) {
	this.resultCh = make(chan readResult, 1024)

	this.wg.Add(1)
	go this.run(ctx, read)
}

func (this *NodeReader) Next() (node *parser.Node, err error) {
	r, ok := <-this.resultCh
	if !ok {
		err = errors.New("NodeReader stopped")
		return
	}
	if r.err != nil {
		err = r.err
		return
	}
	node = &r.node
	return
}
