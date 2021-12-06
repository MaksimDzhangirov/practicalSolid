package application

import (
	"errors"
	"fmt"
	"github.com/MaksimDzhangirov/practicalSOLID/dip/goodCode1/domain"
	"testing"
)

type GetByIDFunc func (id uint) (*domain.User, error)

func (f GetByIDFunc) GetByID(id uint) (*domain.User, error) {
	return f(id)
}

func TestEmailService_SendRegistrationEmail(t *testing.T) {
	service := NewEmailService(GetByIDFunc(func(id uint) (*domain.User, error) {
		return nil, errors.New("error")
	}))
	fmt.Println(service)
	//
	// и после этого просто вызываем сервис
}