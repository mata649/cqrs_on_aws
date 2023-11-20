package domain

import (
	"context"
)

type UserRepository interface {
	Close()
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Get(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user User) error
}
