package increasing

import (
	"context"

	"github.com/mata649/cqrs_on_aws/domain/user"
)

type TaskCounterIncreaserService struct {
	repository user.UserRepository
}

func NewTaskCounterIncreaserService(repository user.UserRepository) TaskCounterIncreaserService {
	return TaskCounterIncreaserService{repository: repository}
}

func (s TaskCounterIncreaserService) IncreaseTasksCreatedCounter(ctx context.Context, userID string) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	user.IncreaseTasksCreatedCounter()
	err = s.repository.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
