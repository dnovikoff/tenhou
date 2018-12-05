package utils

import (
	"sync"
	"sync/atomic"
)

type Routines struct {
	wg        sync.WaitGroup
	lastError atomic.Value
}

func (r *Routines) Start(f func() error) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		err := f()
		if err != nil {
			r.lastError.Store(err)
		}
	}()
}

func (r *Routines) Error() error {
	x := r.lastError.Load()
	if x == nil {
		return nil
	}
	return x.(error)
}

func (r *Routines) Wait() error {
	r.wg.Wait()
	return r.Error()
}
