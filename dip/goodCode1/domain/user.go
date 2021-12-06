package domain

// уровень предметной области

type User struct {
	// какие-то поля
}

type UserRepository interface {
	GetByID(id uint) (*User, error)
}
