package util

import (
	"sync"
	"sync/atomic"
)

type AtomicWaitGroup[T any] struct {
	data *atomic.Pointer[T]
	wg   *sync.WaitGroup
}

// NewAtomicWaitGroup adds one to wg. Store must be called exactly once on the
// returned instance, otherwise any call to Load blocks forever.
func NewAtomicWaitGroup[T any](wg *sync.WaitGroup) *AtomicWaitGroup[T] {
	wg.Add(1)
	return &AtomicWaitGroup[T]{
		data: &atomic.Pointer[T]{},
		wg:   wg,
	}
}

func (a *AtomicWaitGroup[T]) Store(d *T) {
	a.data.Store(d)
	a.wg.Done()
}

func (a *AtomicWaitGroup[T]) Load() *T {
	a.wg.Wait()
	return a.data.Load()
}
