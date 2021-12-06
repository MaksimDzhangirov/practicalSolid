package badCode2

// слишком сильное разбиение
type UserWithFirstName interface {
	FirstName() string
}

type UserWithLastName interface {
	LastName() string
}

type UserWithFullName interface {
	FullName() string
}

// оптимальное разбиение
type UserWithName interface {
	FirstName() string
	LastName() string
	FullName() string
}
