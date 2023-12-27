package inmemory

import (
	"context"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/command"
	"github.com/mata649/cqrs_on_aws/kit/response"
)

type CommandBus struct {
	handlers map[command.Type]command.Handler
}

func NewCommandBus() CommandBus {
	return CommandBus{
		handlers: make(map[command.Type]command.Handler),
	}
}
func (c CommandBus) Dispatch(ctx context.Context, command command.Command) response.Response {
	handler, ok := c.handlers[command.Type()]
	if !ok {
		return response.NewResponse(http.StatusInternalServerError, "invalid command")
	}
	return handler.Handle(ctx, command)
}

func (c CommandBus) Register(commandType command.Type, commandHandler command.Handler) {
	c.handlers[commandType] = commandHandler
}
