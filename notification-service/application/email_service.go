package application

import (
	"fmt"
	"log"
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
	return &EmailService{
		fromEmail: fromEmail,
	}
}

func (es *EmailService) Send(message EmailMessage) error {
	// Mock implementation - just log the email
	log.Printf("Mock Email Sent\n")
	log.Printf("  From: %s\n", es.fromEmail)
	log.Printf("  To: %s\n", message.To)
	log.Printf("  Subject: %s\n", message.Subject)
	log.Printf("  Body: %s\n", message.Body)

	fmt.Printf("Email successfully sent to %s\n", message.To)
	return nil
}
