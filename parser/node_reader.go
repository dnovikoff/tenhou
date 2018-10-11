package parser

import (
	"context"
	"errors"
	"sync"
)

type readResult struct {
	node Node
	err  error
}

type NodeReader struct {
	resultCh     chan readResult
	ReadCallback func(context.Context) (string, error)
}

func NewNodeReader() *NodeReader {
	ch := make(chan readResult)
	// Create with closed channel for correct errors
	close(ch)
	return &NodeReader{resultCh: ch}

}

func (r *NodeReader) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		close(r.resultCh)
	}()
	for {
		message, err := r.ReadCallback(ctx)
		if err != nil {
			r.resultCh <- readResult{err: err}
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
		nodes, err := ParseXML(message)
		if err != nil {
			r.resultCh <- readResult{err: err}
			return
		}
		for _, v := range nodes {
			r.resultCh <- readResult{v, nil}
		}
	}
}

func (r *NodeReader) Start(ctx context.Context) func() {
	r.resultCh = make(chan readResult, 1024)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		r.run(ctx)
		wg.Done()
	}()
	return wg.Wait
}

func (r *NodeReader) Next() (node *Node, err error) {
	res, ok := <-r.resultCh
	if !ok {
		err = errors.New("NodeReader stopped")
		return
	}
	if res.err != nil {
		err = res.err
		return
	}
	node = &res.node
	return
}
