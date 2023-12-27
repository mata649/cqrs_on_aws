package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mata649/cqrs_on_aws/kit/command"
	"github.com/mata649/cqrs_on_aws/kit/query"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/handler"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/middleware"
)

type Server struct {
	commandBus command.Bus
	queryBus   query.Bus
	engine     *chi.Mux
}

func NewServer(commandBus command.Bus, queryBus query.Bus) Server {

	return Server{
		commandBus: commandBus,
		queryBus:   queryBus,
		engine:     chi.NewRouter(),
	}

}

func (s Server) Engine() *chi.Mux {
	return s.engine

}
func (s Server) Route() {
	s.engine.Post("/auth", handler.LoginUserHandler(s.queryBus))
	s.engine.Post("/", handler.CreateUserHandler(s.commandBus))
	s.engine.Delete("/{userID}", middleware.ValidateJWTMiddleware(handler.DeleteUserHandler(s.commandBus)))

}
