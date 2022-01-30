package pool

import (
	"sync"
)

type Job func() error

type pool struct {
	size    int
	wg      sync.WaitGroup
	queue   chan Job
	errOnce sync.Once
	err     error
}

// New returns a new pool.
func New(size int) *pool {
	return &pool{
		size:  size,
		queue: make(chan Job),
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
		go p.worker(i)
	}
	p.wg.Wait()

	// Return first seen error if any.
	return p.err
}

func (p *pool) worker(id int) {
	defer p.wg.Done()
	for {
		job, open := <-p.queue
		if !open {
			break
		}
		err := job()
		if err != nil {
			p.errOnce.Do(func() {
				p.err = err
			})
		}
	}
}
