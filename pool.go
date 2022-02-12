package pool

import (
	"context"
	"sync"
)

type Job func() error

type pool struct {
	size    int
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// New returns a new pool.
func New(size int) *pool {
	return &pool{
		size: size,
	}
}

// Run spawns the numbers of workers defined in the
// pool and starts to execute the provided jobs.
func (p *pool) Run(ctx context.Context, jobs ...Job) error {
	ctx, cancel := context.WithCancel(ctx)

	queue := make(chan Job, len(jobs))
	for _, j := range jobs {
		queue <- j
	}
	close(queue)

	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.worker(ctx, cancel, i, queue)
	}
	p.wg.Wait()
	cancel()

	// Return first seen error if any.
	return p.err
}

func (p *pool) worker(ctx context.Context, cancel func(), id int, queue <-chan Job) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, open := <-queue:
			if !open {
				return
			}
			err := job()
			if err != nil {
				p.errOnce.Do(func() {
					p.err = err
					cancel()
				})
				return
			}
		}
	}
}
