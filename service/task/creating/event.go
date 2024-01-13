package creating

import (
	"context"

	"github.com/mata649/cqrs_on_aws/domain/task"
	"github.com/mata649/cqrs_on_aws/service/user/increasing"
)

type IncreaseTasksCreatedCounterOnTaskCreated struct {
	increaserService increasing.TaskCounterIncreaserService
}

func NewIncreaseCoursesCounterOnCourseCreated(increaserService increasing.TaskCounterIncreaserService) IncreaseTasksCreatedCounterOnTaskCreated {
	return IncreaseTasksCreatedCounterOnTaskCreated{increaserService: increaserService}
}

func (s IncreaseTasksCreatedCounterOnTaskCreated) Handle(ctx context.Context, event task.TaskCreatedEvent) error {
	err := s.increaserService.IncreaseTasksCreatedCounter(ctx, event.UserID())
	if err != nil {
		return err
	}
	return nil
}
