package inmemory

import (
	"context"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/query"
	"github.com/mata649/cqrs_on_aws/kit/response"
)

type QueryBus struct {
	handlers map[query.Type]query.Handler
}

func NewQueryBus() QueryBus {
	return QueryBus{
		handlers: make(map[query.Type]query.Handler),
	}
}
func (q QueryBus) Dispatch(ctx context.Context, query query.Query) response.Response {
	handler, ok := q.handlers[query.Type()]
	if !ok {
		return response.NewResponse(http.StatusInternalServerError, "invalid command")
	}
	return handler.Handle(ctx, query)
}

func (q QueryBus) Register(commandType query.Type, commandHandler query.Handler) {
	q.handlers[commandType] = commandHandler
}
