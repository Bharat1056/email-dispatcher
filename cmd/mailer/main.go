package main

import (
	"fmt"
	"log"
	"time"

	"github.com/user/queue/internal/emailqueue"
)

const (
	queueBufferSize = 6
	workerCount     = 4
	sendLatency     = 50 * time.Millisecond
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("mailer error: %v", err)
	}
}

func run() error {
	queue, err := emailqueue.New(
		emailqueue.Config{
			BufferSize:  queueBufferSize,
			WorkerCount: workerCount,
		},
		emailqueue.NewMailer(sendLatency),
	)
	if err != nil {
		return err
	}

	for _, message := range seedMessages() {
		_ = queue.Enqueue(message)
	}

	queue.Shutdown()
	printFailedJobs(queue.FailedJobs())

	return nil
}

func seedMessages() []emailqueue.Message {
	return []emailqueue.Message{
		{To: "alice@example.com", Body: "Welcome to the queue demo."},
		{To: "bob@example.org", Body: "Your weekly digest is ready."},
		{To: "carol@example.net", Body: "Password reset instructions."},
		{To: "daniel@example.co.in", Body: "Invoice attached."},
		{To: "eve+alerts@example.io", Body: "Alert threshold reached."},
		{To: "frank@", Body: "Missing domain should fail."},
		{To: "", Body: "Empty recipient should fail."},
		{To: "grace.example.com", Body: "Missing at-sign should fail."},
		{To: "@example.com", Body: "Missing local part should fail."},
		{To: "henry@example", Body: "Missing dot in domain should fail."},
		{To: "ivy@example.com", Body: "Reminder: stand-up in 10 minutes."},
		{To: "john.doe@example.dev", Body: "Build completed successfully."},
		{To: "kate_smith@example.ai", Body: "Model retraining started."},
		{To: "leo@example.travel", Body: "Your itinerary is confirmed."},
		{To: "mia@example.school", Body: "Parent-teacher meeting update."},
	}
}

func printFailedJobs(failedJobs []emailqueue.FailedJob) {
	fmt.Printf("--- FAILED JOBS REPORT (%d total) ---\n", len(failedJobs))

	for idx, failedJob := range failedJobs {
		fmt.Printf("%d. email=%s  reason=%s\n", idx+1, failedJob.Email, failedJob.Reason)
	}
}
