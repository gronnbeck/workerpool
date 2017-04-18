package workerpool

import "context"

// New creates a new WorkerPool
func New(noWorkers int) WorkerPool {
	return WorkerPool{
		noWorkers: noWorkers,
		in:        make(chan Work),
		exit:      make(chan bool),
	}
}

// NewBuffered creates a WorkerPool with a channel buffer
func NewBuffered(noWorkers int, bufSize int) WorkerPool {
	return WorkerPool{
		noWorkers: noWorkers,
		in:        make(chan Work, bufSize),
		exit:      make(chan bool),
	}
}

// WorkerPool holds the info needed to run a worker pool
type WorkerPool struct {
	noWorkers int
	in        chan Work
	exit      chan bool
}

// Work is the abstraction of a work
type Work func()

// Start starts the workers in a worker pool
func (pool WorkerPool) Start() {
	for i := 0; i < pool.noWorkers; i++ {
		go func() {
			for {
				select {
				case work := <-pool.in:
					work()
				case <-pool.exit:
					return
				}
			}
		}()
	}
}

// Stop stops a pool of workers
func (pool WorkerPool) Stop() {
	pool.exit <- true
}

// Queue queues work to be completed by WorkerPool
func (pool WorkerPool) Queue(ctx context.Context, work Work) error {
	select {
	case pool.in <- work:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Len returns the length of the worker pool
func (pool WorkerPool) Len() int {
	return len(pool.in)
}
