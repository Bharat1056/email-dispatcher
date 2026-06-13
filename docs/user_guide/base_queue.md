# User Guide: Implementing Custom Queues with BaseQueue

`BaseQueue` is the core foundation for all queues in this project. It provides the essential lifecycle management, thread-safe state tracking, and shutdown signaling needed to build polymorphic, concurrent queues.

If you want to implement a new type of queue in the future (such as a database-backed queue, a priority queue, or a network-backed queue), you should use `BaseQueue` to manage its lifecycle.

---

## What BaseQueue Provides

Under the hood, `BaseQueue` manages:
1. **Queue Identity**: A standard string name for logging and tracking.
2. **Lifecycle State**: Thread-safe tracking of whether the queue is active or shut down.
3. **Shutdown Signal**: A read-only channel (`ShutdownChan`) that closes when the queue is shut down, allowing concurrent workers to wake up and terminate cleanly.

---

## How to Implement a Custom Queue

To create a new queue implementation using `BaseQueue`, follow these steps:

### Step 1: Embed BaseQueue in Your Struct

Always embed `basequeue.BaseQueue` directly in your new queue struct. This gives your struct access to all of `BaseQueue`'s lifecycle methods automatically.

```go
package myqueue

import (
	"sync"
	"github.com/user/queue/internal/basequeue"
)

type CustomQueue struct {
	basequeue.BaseQueue
	mu   sync.Mutex
	data []interface{}
}
```

### Step 2: Define a Constructor

Initialize the embedded `BaseQueue` inside your constructor by calling `basequeue.NewBaseQueue(name)`.

```go
func NewCustomQueue(name string) *CustomQueue {
	return &CustomQueue{
		BaseQueue: basequeue.NewBaseQueue(name),
		data:      make([]interface{}, 0),
	}
}
```

### Step 3: Implement thread-safe Push (Writing to the Queue)

When pushing items, protect the queue state with a mutex. Always check if the queue is closed using `IsClosed()` before writing or modifying data.

```go
func (q *CustomQueue) Push(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 1. Guard against writes to a closed queue
	if q.IsClosed() {
		return fmt.Errorf("cannot push: queue %s is closed", q.Name())
	}

	q.data = append(q.data, item)
	return nil
}
```

### Step 4: Implement non-blocking/blocking Pop (Reading from the Queue)

In your consumer or worker loops, listen to `ShutdownChan()` using a `select` statement. This ensures your workers do not block indefinitely when the queue is shut down.

```go
func (q *CustomQueue) Pop() (interface{}, error) {
	select {
	case <-q.ShutdownChan():
		// 1. Wake up and exit immediately when shutdown is signaled
		return nil, fmt.Errorf("queue %s is closed", q.Name())
	default:
		// 2. Perform non-blocking read or check local buffer
		q.mu.Lock()
		defer q.mu.Unlock()
		if len(q.data) == 0 {
			return nil, fmt.Errorf("queue is empty")
		}
		item := q.data[0]
		q.data = q.data[1:]
		return item, nil
	}
}
```

### Step 5: Implement the Shutdown Method

Your custom queue's `Shutdown()` method must:
1. Acquire your queue's lock.
2. Check if it's already closed using `IsClosed()`.
3. Call the parent `BaseQueue.Shutdown()`, which closes the shutdown channel.
4. Clean up any custom local resources (e.g. closing database connections or channels).

```go
func (q *CustomQueue) Shutdown() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsClosed() {
		return
	}

	// 1. Trigger the shutdown channel and set closed flag
	q.BaseQueue.Shutdown()

	// 2. Perform local cleanups (e.g., clear buffer)
	q.data = nil
}
```
