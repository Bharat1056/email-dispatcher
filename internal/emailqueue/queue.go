package emailqueue

import (
	"fmt"
	"sync"
)

// Message is the unit processed by the worker pool.
type Message struct {
	To   string
	Body string
}

// FailedJob captures queue-time and delivery-time failures.
type FailedJob struct {
	Email  string
	Reason string
}

// Sender defines the delivery contract so SMTP can replace the simulator later.
type Sender interface {
	Send(Message) error
}

// Config controls queue capacity and worker fan-out.
type Config struct {
	BufferSize  int
	WorkerCount int
}

// Queue owns the message channel, workers, and failure tracking.
type Queue struct {
	messages chan Message
	sender   Sender

	wg sync.WaitGroup

	mu     sync.Mutex
	failed []FailedJob
}

// New constructs a queue and starts its workers immediately.
func New(cfg Config, sender Sender) (*Queue, error) {
	if cfg.BufferSize <= 0 {
		return nil, fmt.Errorf("buffer size must be greater than zero")
	}

	if cfg.WorkerCount <= 0 {
		return nil, fmt.Errorf("worker count must be greater than zero")
	}

	if sender == nil {
		return nil, fmt.Errorf("sender is required")
	}

	q := &Queue{
		messages: make(chan Message, cfg.BufferSize),
		sender:   sender,
		failed:   make([]FailedJob, 0),
	}

	q.startWorkers(cfg.WorkerCount)

	return q, nil
}

// Enqueue validates a message before allowing it into the channel.
func (q *Queue) Enqueue(message Message) error {
	if err := validateEmail(message.To); err != nil {
		q.recordFailure(message.To, err.Error())
		fmt.Printf("[queue] invalid | to=%s | reason=%s\n", message.To, err.Error())
		return err
	}

	q.messages <- message
	return nil
}

// Shutdown closes the channel and waits for all workers to drain it.
func (q *Queue) Shutdown() {
	close(q.messages)
	q.wg.Wait()
}

// FailedJobs returns a snapshot of recorded failures.
func (q *Queue) FailedJobs() []FailedJob {
	q.mu.Lock()
	defer q.mu.Unlock()

	failed := make([]FailedJob, len(q.failed))
	copy(failed, q.failed)

	return failed
}

func (q *Queue) startWorkers(count int) {
	for workerID := 1; workerID <= count; workerID++ {
		q.wg.Add(1)

		go func(id int) {
			defer q.wg.Done()

			fmt.Printf("[worker-%d] started\n", id)

			for message := range q.messages {
				fmt.Printf("[worker-%d] processing | to=%s | body=%s\n", id, message.To, message.Body)

				if err := q.sender.Send(message); err != nil {
					q.recordFailure(message.To, err.Error())
					fmt.Printf("[worker-%d] delivery failed | to=%s | reason=%s\n", id, message.To, err.Error())
				}
			}

			fmt.Printf("[worker-%d] exiting\n", id)
		}(workerID)
	}
}

func (q *Queue) recordFailure(email, reason string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.failed = append(q.failed, FailedJob{
		Email:  email,
		Reason: reason,
	})
}
