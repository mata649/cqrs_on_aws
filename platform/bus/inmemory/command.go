package inmemory

import (
	"context"
	"net/http"

	"github.com/mata649/cqrs_on_aws/domain/command"
	"github.com/mata649/cqrs_on_aws/domain/response"
)

type CommandBus struct {
	handlers map[command.Type]command.Handler
}

func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[command.Type]command.Handler),
	}
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd command.Command) response.Response {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return response.NewResponse(http.StatusInternalServerError, "Command not found")
	}
	return handler.Handle(ctx, cmd)
}

func (b *CommandBus) Register(cmdType command.Type, handler command.Handler) {
	b.handlers[cmdType] = handler
}
