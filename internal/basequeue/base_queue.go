package basequeue

import (
	"sync"
)

// BaseQueue defines the core identity and lifecycle synchronization for polymorphic queues.
type BaseQueue struct {
	name         string
	mu           sync.Mutex
	closed       bool
	shutdownChan chan struct{}
}

// NewBaseQueue constructs a new BaseQueue structure.
func NewBaseQueue(name string) BaseQueue {
	return BaseQueue{
		name:         name,
		shutdownChan: make(chan struct{}),
	}
}

func (bq *BaseQueue) init() {
	if bq.shutdownChan == nil {
		bq.shutdownChan = make(chan struct{})
	}
}

// Name returns the queue name.
func (bq *BaseQueue) Name() string {
	return bq.name
}

// IsClosed returns a thread-safe check of whether the queue is shut down.
func (bq *BaseQueue) IsClosed() bool {
	bq.mu.Lock()
	defer bq.mu.Unlock()
	return bq.closed
}

// Shutdown transitions the queue into a closed state and closes the shutdown channel thread-safely.
func (bq *BaseQueue) Shutdown() {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	bq.init()
	if bq.closed {
		return
	}

	bq.closed = true
	close(bq.shutdownChan)
}

// ShutdownChan returns the read-only shutdown channel for queue components to await shutdown signals.
func (bq *BaseQueue) ShutdownChan() <-chan struct{} {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	bq.init()
	return bq.shutdownChan
}
