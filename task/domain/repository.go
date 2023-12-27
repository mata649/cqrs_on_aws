package domain

import "context"

type TaskRepository interface {
	Close()
	Create(ctx context.Context, task Task) error
	GetByID(ctx context.Context, id string) (Task, error)
	Get(ctx context.Context) ([]Task, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, task Task) error
}
