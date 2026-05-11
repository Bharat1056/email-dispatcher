package emailqueue

import (
	"fmt"
	"time"
)

// Mailer simulates message delivery behind the Sender interface.
type Mailer struct {
	latency time.Duration
}

// NewMailer returns a sender implementation that can be swapped later.
func NewMailer(latency time.Duration) *Mailer {
	return &Mailer{latency: latency}
}

// Send simulates network latency and logs successful delivery.
func (m *Mailer) Send(message Message) error {
	time.Sleep(m.latency)
	fmt.Printf("[mailer] sent | to=%s | body=%s\n", message.To, message.Body)
	return nil
}
