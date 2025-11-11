package sender

import (
	"fmt"

	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridSender struct {
	config config.SendGridConfig
	client *sendgrid.Client
}

func NewSendGridSender(cfg config.SendGridConfig) (*SendGridSender, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("SendGrid API key is required")
	}

	client := sendgrid.NewSendClient(cfg.APIKey)

	return &SendGridSender{
		config: cfg,
		client: client,
	}, nil
}

func (s *SendGridSender) Send(to, subject, body string) error {
	from := mail.NewEmail("Notification System", "noreply@example.com")
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, "", body)

	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid error: status %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}

func (s *SendGridSender) GetProviderName() string {
	return "sendgrid"
}
