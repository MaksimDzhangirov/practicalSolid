package badCode

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	Amount     int       `gorm:"column:amount" json:"amount" validate:"required"`
	CurrencyID int       `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Time       time.Time `gorm:"column:time" json:"time" validate:"required"`
}
