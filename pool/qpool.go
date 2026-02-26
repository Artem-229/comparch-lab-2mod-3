package pool

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrStopped = errors.New("pool was stopped")
)

type ProcessRoutine[T any] func(context.Context, int, T)

type QPool[T any] struct {
	tasks   chan T
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mtx     sync.Mutex
	stopped bool
}

func NewQPool[T any](workers int, routine ProcessRoutine[T]) *QPool[T] {
	p := &QPool[T]{
		tasks: make(chan T, workers),
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.wg.Add(workers)

	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer p.wg.Done()
			for task := range p.tasks {
				routine(p.ctx, i, task)
			}
		}()
	}

	return p
}

func (p *QPool[T]) Push(task T) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.stopped {
		return ErrStopped
	}

	select {
	case <-p.ctx.Done():
		return ErrStopped
	case p.tasks <- task:
		return nil
	}
}

func (p *QPool[T]) Stop() {
	p.cancel()

	p.mtx.Lock()
	if !p.stopped {
		p.stopped = true
		close(p.tasks)
	}
	p.mtx.Unlock()

	p.wg.Wait()
}
