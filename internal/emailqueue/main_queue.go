package emailqueue

import (
	"errors"
	"sync"

	"github.com/user/queue/internal/basequeue"
)

// MainQueue embeds BaseQueue and uses a Go channel for worker handoff.
type MainQueue struct {
	basequeue.BaseQueue
	mu sync.Mutex // MainQueue's own mutex protecting the channel
	ch chan *Job
}

// NewMainQueue initializes a new MainQueue instance.
func NewMainQueue(name string, size int) *MainQueue {
	return &MainQueue{
		BaseQueue: basequeue.NewBaseQueue(name),
		ch:        make(chan *Job, size),
	}
}

// Push adds a job to the in-memory queue.
func (mq *MainQueue) Push(job *Job) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if mq.IsClosed() {
		return errors.New("cannot push: main queue is closed")
	}

	mq.ch <- job
	return nil
}

// Pop retrieves the next job, blocking until one is available or the queue is closed.
func (mq *MainQueue) Pop() (*Job, error) {
	select {
	case <-mq.ShutdownChan():
		return nil, errors.New("main queue is closed")
	case job, ok := <-mq.ch:
		if !ok {
			return nil, errors.New("main queue channel closed")
		}
		return job, nil
	}
}

// Size returns the current number of items buffered in memory.
func (mq *MainQueue) Size() int {
	return len(mq.ch)
}

// Shutdown transitions the queue to a closed state, closes the job channel, and notifies all waiting workers.
func (mq *MainQueue) Shutdown() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if mq.IsClosed() { // Guard against double-closing the channel
		return
	}
	mq.BaseQueue.Shutdown()
	close(mq.ch)
}

// Channel exposes the raw channel for draining purposes.
func (mq *MainQueue) Channel() chan *Job {
	return mq.ch
}
