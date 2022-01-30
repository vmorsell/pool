package pool

import (
	"context"
	"sync"
)

type Job func() error

type pool struct {
	ctx    context.Context
	cancel func()

	size    int
	wg      sync.WaitGroup
	queue   chan Job
	errOnce sync.Once
	err     error
}

// New returns a new pool.
func New(size int) *pool {
	_, pool := NewWithContext(context.Background(), size)
	return pool
}

// NewWithContext returns a new context-aware pool.
func NewWithContext(ctx context.Context, size int) (context.Context, *pool) {
	ctx, cancel := context.WithCancel(ctx)
	return ctx, &pool{
		ctx:    ctx,
		cancel: cancel,
		size:   size,
		queue:  make(chan Job),
	}
}

// Run spawns the numbers of workers defined in the
// pool and starts to execute the provided jobs.
func (p *pool) Run(jobs ...Job) error {
	go func() {
		for _, j := range jobs {
			p.queue <- j
		}
		close(p.queue)
	}()

	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	p.wg.Wait()

	// Return first seen error if any.
	return p.err
}

func (p *pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case job, open := <-p.queue:
			if !open {
				return
			}
			err := job()
			if err != nil {
				p.errOnce.Do(func() {
					p.err = err
					p.cancel()
				})
			}
		}
	}
}
