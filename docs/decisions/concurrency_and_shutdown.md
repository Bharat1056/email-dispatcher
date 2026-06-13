# Design Decisions: Concurrency and Shutdown Architecture

This document details the architectural decisions made for the synchronization and shutdown mechanisms of `BaseQueue` and `MainQueue`.

---

## 1. Why `Push` uses `IsClosed()` but `Pop` does not

### The Challenge
In Go, channel operations behave differently depending on the operation:
* **Sending** (`mq.ch <- job`) to a closed channel causes a **runtime panic**.
* **Receiving** (`<-mq.ch`) from a closed channel is **safe** and returns the zero value along with `ok = false`.

### The Decision
* To prevent runtime panics, `Push` must be guaranteed never to send to a closed channel. We achieve this by acquiring a mutex lock (`mq.mu`), checking `IsClosed()`, and only writing if the queue is open.
* Since receiving from a closed channel never panics, `Pop` does not need to guard against closed channels with a lock. Instead, it handles channel closures gracefully via the second return value `ok`:
  ```go
  case job, ok := <-mq.ch:
      if !ok {
          return nil, errors.New("main queue channel closed")
      }
  ```

---

## 2. Why `Pop` uses a Select statement instead of Mutex Checks

### The Challenge
`Pop` is a blocking operation. If there are no jobs in the queue, workers must sleep until a job arrives. 
* If `Pop` held a mutex lock while waiting, it would block all other goroutines from pushing jobs or calling shutdown.
* If `Pop` checked `IsClosed()` without holding a lock before blocking, a race condition could occur: the queue could shut down right after the check, leaving the worker blocked indefinitely on `<-mq.ch`.

### The Decision
We use a `select` statement that listens on both the shutdown channel and the job channel:
```go
select {
case <-mq.ShutdownChan():
    return nil, errors.New("main queue is closed")
case job, ok := <-mq.ch:
    // ...
}
```
This allows workers to sleep without holding any locks. When `Shutdown()` is called and the shutdown channel is closed, all blocked workers wake up immediately and exit cleanly.

---

## 3. Handling Pseudo-Random Select Choices at Shutdown

### The Challenge
When a queue is shut down, both `mq.ch` and the shutdown channel are closed. In Go, if multiple cases in a `select` block are ready at the same time, the runtime chooses one **pseudo-randomly**.

### The Decision
If the worker selects the job channel case during shutdown:
1. **If jobs are still in the buffer**: It retrieves the job and processes it. This is a feature, as we want to drain remaining jobs before stopping.
2. **If the buffer is empty**: The read returns `ok = false`, and the worker returns a "channel closed" error.

Both cases are fully safe, and no worker gets stuck.

---

## 4. Preventing Double-Close Panics

### The Challenge
In Go, calling `close()` twice on the same channel causes a panic. If `Shutdown()` is called multiple times, we must ensure we only close `mq.ch` once.

### The Decision
We guard the shutdown state using `mq.IsClosed()` inside `Shutdown()` under a mutex lock:
```go
mq.mu.Lock()
defer mq.mu.Unlock()

if mq.IsClosed() {
    return // Prevents running the shutdown code again
}
mq.BaseQueue.Shutdown()
close(mq.ch) // Safe to close exactly once
```
Because the lock prevents concurrent executions of `Shutdown()`, the channel is guaranteed to be closed exactly once.

---

## 5. Security and Memory Safety

### The Challenge
Is there any risk of data exposure or hacking during the shutdown sequence?

### The Decision
* **Memory Isolation**: Go channels and jobs exist entirely in the server's private RAM. They are not exposed to the network, meaning external attackers cannot access or "sniff" the data.
* **Strict Concurrency Guards**: All state mutations are serialized through Go mutexes. There is no race condition where a writer can push data to a partially torn-down queue.
