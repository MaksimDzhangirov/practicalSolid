package badCode1

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID
	//
	// какие-то поля
	//
}

type UserRepository interface {
	Update(ctx context.Context, user User) error
}

type DBUserRepository struct {
	db *gorm.DB
}

func (r *DBUserRepository) Update(ctx context.Context, user User) error {
	return r.db.WithContext(ctx).Delete(user).Error
}