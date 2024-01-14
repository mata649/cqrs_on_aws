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

func NewCreateTaskService(repository task.TaskRepository, eventBus event.Bus) CreateTaskService {
	return CreateTaskService{repository: repository, eventBus: eventBus}
}

func (s CreateTaskService) Create(ctx context.Context, id, title, description, userID string, createdAt time.Time) response.Response {

	task, err := task.NewTask(id, title, description, createdAt, userID)
	if err != nil {
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}

	err = s.repository.Create(ctx, task)

	if err != nil {
		log.Println("Error creating task: ", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	err = s.eventBus.Publish(ctx, task.PullEvents())
	if err != nil {
		log.Println("Error publishing events: ", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	return response.NewResponse(http.StatusCreated, "Task created successfully")

}
