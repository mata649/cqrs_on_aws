package creating

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/domain/event"
	"github.com/mata649/cqrs_on_aws/domain/response"
	"github.com/mata649/cqrs_on_aws/domain/task"
)

type CreateTaskService struct {
	repository task.TaskRepository
	eventBus   event.Bus
}

func NewCreateTaskService(repository task.TaskRepository) CreateTaskService {
	return CreateTaskService{repository: repository}
}

func (s CreateTaskService) Create(ctx context.Context, id, title, description, userID string, createdAt time.Time) response.Response {
	log.Println(userID)
	task, err := task.NewTask(id, title, description, createdAt, userID)
	if err != nil {
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}

	err = s.repository.Create(ctx, task)
	log.Println(err)
	if err != nil {
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	// s.eventBus.Publish(ctx, task.PullEvents())
	return response.NewResponse(http.StatusCreated, "Task created successfully")

}
