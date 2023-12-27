package creating

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/kit/command"
	"github.com/mata649/cqrs_on_aws/kit/response"
)

const CreateUserCommandType command.Type = "command.creating.user"

type CreateUserCommand struct {
	id       string
	username string
	password string
	createAt time.Time
}

func NewCreateUserCommand(id, username, password string, createAt time.Time) CreateUserCommand {
	return CreateUserCommand{
		id:       id,
		username: username,
		password: password,
		createAt: createAt,
	}
}

func (c CreateUserCommand) Type() command.Type {
	return CreateUserCommandType
}

type CreateUserCommandHandler struct {
	service CreateUserService
}

func NewCreateUserCommandHandler(service CreateUserService) CreateUserCommandHandler {
	return CreateUserCommandHandler{service: service}
}

func (h CreateUserCommandHandler) Handle(ctx context.Context, command command.Command) response.Response {
	courseCommand, ok := command.(CreateUserCommand)
	if !ok {
		log.Printf("Invalid command type: %T\n", command)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	return h.service.Create(ctx, courseCommand.id, courseCommand.username, courseCommand.password, courseCommand.createAt)
}
