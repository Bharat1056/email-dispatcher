# User Guide: Using MainQueue

`MainQueue` is an in-memory, thread-safe queue implementation designed for high-performance worker handoff. It embeds `BaseQueue` and utilizes a Go buffered channel to transfer `Job` objects between producers (publishers) and consumers (workers).

---

## How MainQueue Works

`MainQueue` acts as a fast, in-memory buffer. It allows you to:
1. **Push jobs** into the queue asynchronously.
2. **Pop jobs** concurrently, blocking the workers automatically until a job becomes available.
3. **Shutdown safely** without causing panic or losing active workers.

---

## Core API & Usage

### 1. Initializing the Queue
To create a new `MainQueue`, specify the queue name and the buffer size. The buffer size dictates how many jobs the queue can hold in memory before `Push` blocks.

```go
import "github.com/user/queue/internal/emailqueue"

// Initialize a queue named "email-delivery" with a buffer of 100 jobs
mq := emailqueue.NewMainQueue("email-delivery", 100)
```

### 2. Pushing Jobs to the Queue
Producers use `Push` to add jobs to the queue. 
* If the buffer is full, `Push` will block until space becomes available.
* If the queue is closed, `Push` returns an error.

```go
job := &emailqueue.Job{ /* job details */ }

err := mq.Push(job)
if err != nil {
	log.Printf("Failed to push job: %v", err)
}
```

### 3. Popping Jobs from the Queue (Workers)
Workers use `Pop` to retrieve jobs. `Pop` blocks if the queue is empty, and automatically wakes up and returns when:
* A new job is pushed to the queue.
* The queue is shut down (returns a shutdown error).

```go
go func() {
	for {
		job, err := mq.Pop()
		if err != nil {
			log.Printf("Worker stopping: %v", err)
			return
		}
		
		// Process the job
		process(job)
	}
}()
```

### 4. Shutting Down the Queue
To shut down the queue, call `Shutdown()`. This will:
1. Prevent any new jobs from being pushed.
2. Close the in-memory channel.
3. Wake up and cleanly terminate all workers currently waiting inside `Pop()`.

```go
mq.Shutdown()
```
