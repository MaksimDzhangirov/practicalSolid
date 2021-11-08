package badCode

import (
	"errors"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Amount     int `gorm:"column:amount"`
	CurrencyID int `gorm:"column:currency_id"`
}

func (w *Wallet) Withdraw(amount int) error {
	if amount > w.Amount {
		return errors.New("there is no enough money in wallet")
	}

	w.Amount -= amount

	return nil
}