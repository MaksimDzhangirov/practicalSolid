package application

import (
	"fmt"
	"github.com/MaksimDzhangirov/practicalSOLID/dip/badCode1/infrastructure"
)

// уровень прикладных операций

type EmailService struct {
	repository *infrastructure.UserRepository
	// какой-то отправитель электронных писем
}

func NewEmailService(repository *infrastructure.UserRepository) *EmailService {
	return &EmailService{
		repository: repository,
	}
}

func (s *EmailService) SendRegistrationEmail(userID uint) error {
	user, err := s.repository.GetByID(userID)
	if err != nil {
		return err
	}

	fmt.Println(user)
	// отправляем электронное письмо
	return nil
}
