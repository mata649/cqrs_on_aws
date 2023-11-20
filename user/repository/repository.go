package repository

import (
	"context"

	"github.com/mata649/cqrs_on_aws/user/domain"
)

type UserRepository interface {
	Close()
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id string) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	Get(ctx context.Context) ([]domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user domain.User) error
}

var userRepository UserRepository

func SetUserRepository(repo UserRepository) {
	userRepository = repo
}

func Close() {
	userRepository.Close()
}
func Create(ctx context.Context, user domain.User) error {
	return userRepository.Create(ctx, user)
}
func GetByID(ctx context.Context, id string) (domain.User, error) {
	return userRepository.GetByID(ctx, id)
}
func GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return userRepository.GetByUsername(ctx, username)
}
func Get(ctx context.Context) ([]domain.User, error) {
	return userRepository.Get(ctx)
}
func Delete(ctx context.Context, id string) error {
	return userRepository.Delete(ctx, id)
}

func Update(ctx context.Context, user domain.User) error {
	return userRepository.Update(ctx, user)
}
