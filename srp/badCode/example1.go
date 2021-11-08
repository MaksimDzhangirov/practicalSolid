package badCode

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/smtp"
)

type EmailGorm struct {
	gorm.Model
	From    string
	To      string
	Subject string
	Message string
}

type EmailService struct {
	db           *gorm.DB
	smtpHost     string
	smtpPassword string
	smtpPort     int
}

func NewEmailService(db *gorm.DB, smtpHost string, smtpPassword string, smtpPort int) *EmailService {
	return &EmailService{
		db:           db,
		smtpHost:     smtpHost,
		smtpPassword: smtpPassword,
		smtpPort:     smtpPort,
	}
}

func (s *EmailService) Send(from string, to string, subject string, message string) error {
	email := EmailGorm{
		From:    from,
		To:      to,
		Subject: subject,
		Message: message,
	}

	err := s.db.Create(&email).Error
	if err != nil {
		log.Println(err)
		return err
	}

	auth := smtp.PlainAuth("", from, s.smtpPassword, s.smtpHost)

	server := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

	err = smtp.SendMail(server, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
