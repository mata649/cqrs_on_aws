package creating

import (
	"context"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/domain/command"
	"github.com/mata649/cqrs_on_aws/domain/response"
)

const CreateTaskCommandType command.Type = "command.task.create"

type CreateTaskCommand struct {
	id          string
	title       string
	description string
	createdAt   time.Time
	userID      string
}

func NewCreateTaskCommand(id, title, description, userID string, createdAt time.Time) CreateTaskCommand {
	return CreateTaskCommand{
		id:          id,
		title:       title,
		description: description,
		createdAt:   createdAt,
		userID:      userID,
	}
}

func (c CreateTaskCommand) Type() command.Type {
	return CreateTaskCommandType
}

type CreateTaskCommandHandler struct {
	service CreateTaskService
}

func NewCreateTaskCommandHandler(service CreateTaskService) CreateTaskCommandHandler {
	return CreateTaskCommandHandler{
		service: service,
	}
}

func (h CreateTaskCommandHandler) Handle(ctx context.Context, cmd command.Command) response.Response {
	createTaskCommand, ok := cmd.(CreateTaskCommand)
	if !ok {
		return response.NewResponse(http.StatusInternalServerError, "Invalid command")
	}
	return h.service.Create(ctx, createTaskCommand.id, createTaskCommand.title, createTaskCommand.description, createTaskCommand.userID, createTaskCommand.createdAt)
}
