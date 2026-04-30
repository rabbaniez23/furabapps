package mock

import "context"

// MockEmailSender is a test double for outbound email sender.
type MockEmailSender struct {
	SendFn    func(ctx context.Context, receiverEmail, subject, body string) error
	SendCall  int
	LastEmail string
	LastSubj  string
	LastBody  string
}

// Send records calls and delegates behavior to SendFn.
func (m *MockEmailSender) Send(ctx context.Context, receiverEmail, subject, body string) error {
	m.SendCall++
	m.LastEmail = receiverEmail
	m.LastSubj = subject
	m.LastBody = body
	if m.SendFn != nil {
		return m.SendFn(ctx, receiverEmail, subject, body)
	}
	return nil
}
