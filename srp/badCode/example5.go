package badCode

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	db        *gorm.DB
	Username  string
	Firstname string
	Lastname  string
	Birthday  time.Time
	//
	// какие-то другие поля
	//
}

func (u User) IsAdult() bool {
	return u.Birthday.AddDate(18, 0, 0).Before(time.Now())
}

func (u User) Save() error {
	return u.db.Exec("INSERT INTO users ...", u.Username, u.Firstname, u.Lastname, u.Birthday).Error
}