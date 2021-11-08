package badCode2

import (
	"context"
	"github.com/MaksimDzhangirov/practicalSOLID/lsp/badCode1"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user badCode1.User) (*badCode1.User, error)
	Update(ctx context.Context, user badCode1.User) error
}

type DBUserRepository struct {
	db *gorm.DB
}

func (r *DBUserRepository) Create(ctx context.Context, user badCode1.User) (*badCode1.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	return &user, err
}

func (r *DBUserRepository) Update(ctx context.Context, user badCode1.User) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

type MemoryUserRepository struct {
	users map[uuid.UUID]badCode1.User
}

func (r *MemoryUserRepository) Create(_ context.Context, user badCode1.User) (*badCode1.User, error) {
	if r.users == nil {
		r.users = map[uuid.UUID]badCode1.User{}
	}
	user.ID = uuid.New()
	r.users[user.ID] = user

	return &user, nil
}

func (r *MemoryUserRepository) Update(_ context.Context, user badCode1.User) error {
	if r.users == nil {
		r.users = map[uuid.UUID]badCode1.User{}
	}
	r.users[user.ID] = user

	return nil
}