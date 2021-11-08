package goodCode

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailSender interface {
	Send(from string, to string, subject string, message string) error
}

type EmailSMTPSender struct {
	smtpHost     string
	smtpPassword string
	smtpPort     int
}

func NewEmailSender(smtpHost string, smtpPassword string, smtpPort int) EmailSender {
	return &EmailSMTPSender{
		smtpHost:     smtpHost,
		smtpPassword: smtpPassword,
		smtpPort:     smtpPort,
	}
}

func (s *EmailSMTPSender) Send(from string, to string, subject string, message string) error {
	auth := smtp.PlainAuth("", from, s.smtpPassword, s.smtpHost)

	server := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	err := smtp.SendMail(server, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
