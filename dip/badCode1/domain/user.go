package domain

// уровень предметной области

type User struct {
	ID uint `gorm:"primaryKey;column:id"`
	// какие-то поля
}

type UserRepository interface {
	GetByID(id uint) (*User, error)
}
