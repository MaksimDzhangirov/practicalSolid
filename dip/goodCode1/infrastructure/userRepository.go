package infrastructure

import (
	"github.com/MaksimDzhangirov/practicalSOLID/dip/goodCode1/domain"
	"gorm.io/gorm"
)

// инфраструктурный уровень

type UserGorm struct {
	// какие-то поля
}

func (g UserGorm) ToUser() *domain.User {
	return &domain.User{
		// какие-то поля
	}
}

type UserDatabaseRepository struct {
	db *gorm.DB
}

var _ domain.UserRepository = &UserDatabaseRepository{}

/*
type UserRedisRepository struct {

}

type UserCassandraRepository struct {

}
*/

func NewUserDatabaseRepository(db *gorm.DB) domain.UserRepository {
	return &UserDatabaseRepository{
		db: db,
	}
}

func (r *UserDatabaseRepository) GetByID(id uint) (*domain.User, error) {
	user := UserGorm{}
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user.ToUser(), nil
}
