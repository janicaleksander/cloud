package application

import (
	"log/slog"
)

type EmailService struct {
	fromEmail string
}

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

func NewEmailService(fromEmail string) *EmailService {
	slog.Info("Initializing EmailService", "fromEmail", fromEmail)
	return &EmailService{
		fromEmail: fromEmail,
	}
}

func (es *EmailService) Send(message EmailMessage) error {
	slog.Info("Sending email", "to", message.To, "subject", message.Subject)
	slog.Info("Email content", "body", message.Body)
	slog.Info("Email sent successfully", "to", message.To)
	return nil
}
