package query

import (
	"context"

	"github.com/mata649/cqrs_on_aws/domain/response"
)

type Type string

type Bus interface {
	Dispatch(context.Context, Query) response.Response
	Register(Type, Handler)
}

type Query interface {
	Type() Type
}

type Handler interface {
	Handle(context.Context, Query) response.Response
}
