package sender

// EmailSender interface for different email providers
type EmailSender interface {
	Send(to, subject, body string) error
	GetProviderName() string
}
