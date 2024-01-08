package command

import (
	"context"

	"github.com/mata649/cqrs_on_aws/domain/response"
)

type Bus interface {
	Dispatch(context.Context, Command) response.Response
	Register(Type, Handler)
}

//go:generate mockery --case=snake --outpkg=commandmocks --output=commandmocks --name=Bus
type Type string

type Command interface {
	Type() Type
}

type Handler interface {
	Handle(context.Context, Command) response.Response
}
