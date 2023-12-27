package deleting

import (
	"context"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/command"
	"github.com/mata649/cqrs_on_aws/kit/response"
)

const DeleteUserCommandType command.Type = "command.deleting.user"

type DeleteUserCommand struct {
	userID        string
	currentUserID string
}

func NewDeleteUserCommand(userID, currentUserID string) DeleteUserCommand {
	return DeleteUserCommand{
		userID:        userID,
		currentUserID: currentUserID,
	}
}

func (c DeleteUserCommand) Type() command.Type {
	return DeleteUserCommandType
}

type DeleteUserCommandHandler struct {
	service DeleteUserService
}

func NewDeleteUserCommandHandler(service DeleteUserService) DeleteUserCommandHandler {
	return DeleteUserCommandHandler{service: service}
}

func (h DeleteUserCommandHandler) Handle(ctx context.Context, command command.Command) response.Response {
	userCommand, ok := command.(DeleteUserCommand)
	if !ok {
		return response.NewResponse(http.StatusInternalServerError, "invalid command")
	}
	return h.service.Delete(ctx, userCommand.userID, userCommand.currentUserID)
}
