package task

import (
	"github.com/go-chi/chi/v5"
	"github.com/mata649/cqrs_on_aws/domain/command"
	"github.com/mata649/cqrs_on_aws/platform/server/task/handler"
)

type Server struct {
	commandBus command.Bus

	engine *chi.Mux
}

func NewServer(commandBus command.Bus) Server {

	return Server{
		commandBus: commandBus,
		engine:     chi.NewRouter(),
	}

}

func (s Server) Engine() *chi.Mux {
	return s.engine

}

func (s Server) SetupRoutes() {

	s.engine.Post("/tasks", handler.CreateTaskHandler(s.commandBus))

}
