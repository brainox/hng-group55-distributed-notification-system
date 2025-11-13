package sender

import (
	"fmt"
	"net/smtp"

	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/config"
)

type SMTPSender struct {
	config config.SMTPConfig
	auth   smtp.Auth
}

func NewSMTPSender(cfg config.SMTPConfig) (*SMTPSender, error) {
	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("SMTP username and password are required")
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	return &SMTPSender{
		config: cfg,
		auth:   auth,
	}, nil
}

func (s *SMTPSender) Send(to, subject, body string) error {
	// Build MIME message
	message := s.buildMessage(s.config.Username, to, subject, body)

	// Send email
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, s.auth, s.config.Username, []string{to}, []byte(message))

	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}

func (s *SMTPSender) GetProviderName() string {
	return "smtp"
}

func (s *SMTPSender) buildMessage(from, to, subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	return message
}
