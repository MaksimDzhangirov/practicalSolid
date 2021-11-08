package goodCode1

import (
	"context"
	"github.com/MaksimDzhangirov/practicalSOLID/lsp/badCode1"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user badCode1.User) (*badCode1.User, error)
	Update(ctx context.Context, user badCode1.User) error
}

type MySQLUserRepository struct {
	db *gorm.DB
}

type CassandraUserRepository struct {
	session *gocql.Session
}

type UserCache interface {
	Create(user badCode1.User)
	Update(user badCode1.User)
}

type MemoryUserCache struct {
	users map[uuid.UUID]badCode1.User
}
