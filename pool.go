package pool

import (
	"golang.org/x/sync/errgroup"
)

type Job func() error

type pool struct {
	size  int
	queue chan Job
}

// New returns a new pool.
func New(size int) *pool {
	return &pool{
		size:  size,
		queue: make(chan Job),
	}
}

// Queue runs the provided jobs and returns the first non-nil error if any.
// No more jobs are consumed from the queue if an error occurs.
func (p *pool) Run(jobs ...Job) error {
	go func() {
		for _, j := range jobs {
			p.queue <- j
		}
		close(p.queue)
	}()

	g := new(errgroup.Group)
	for i := 0; i < p.size; i++ {
		g.Go(func() error {
			return worker(p.queue)
		})
	}
	return g.Wait()
}

func worker(jobs <-chan Job) error {
	for job := range jobs {
		if err := job(); err != nil {
			return err
		}
	}
	return nil
}
