package application

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MaksimDzhangirov/practicalSOLID/dip/badCode1/infrastructure"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestEmailService_SendRegistrationEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		DSN:        "dummy",
		DriverName: "mysql",
		Conn:       db,
	})
	finalDB, err := gorm.Open(dialector, &gorm.Config{})

	repository := infrastructure.NewUserRepository(finalDB)
	service := NewEmailService(repository)
	fmt.Println(service, mock)
	//
	// большой фрагмент кода для имитации SQL-запросов
	//
	// а затем сам тест
}
