package goodCode

import (
	"gorm.io/gorm"
	"log"
)

type EmailGorm struct {
	gorm.Model
	From    string
	To      string
	Subject string
	Message string
}

type EmailRepository interface {
	Save(from string, to string, subject string, message string) error
}

type EmailDBRepository struct {
	db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) EmailRepository {
	return &EmailDBRepository{
		db: db,
	}
}

func (r *EmailDBRepository) Save(from string, to string, subject string, message string) error {
	email := EmailGorm{
		From:    from,
		To:      to,
		Subject: subject,
		Message: message,
	}

	err := r.db.Create(&email).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
