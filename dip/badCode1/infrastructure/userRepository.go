package infrastructure

import (
	"github.com/MaksimDzhangirov/practicalSOLID/dip/badCode1/domain"
	"gorm.io/gorm"
)

// инфраструктурный уровень

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	user := domain.User{}
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
