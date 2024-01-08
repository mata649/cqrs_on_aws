package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/mata649/cqrs_on_aws/platform/server/middleware"
	"github.com/mata649/cqrs_on_aws/platform/server/user/handler"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

type Server struct {
	userService service.UserService
	engine      *chi.Mux
}

func NewServer(userService service.UserService) Server {

	return Server{
		userService: userService,

		engine: chi.NewRouter(),
	}

}

func (s Server) Engine() *chi.Mux {
	return s.engine

}

func (s Server) SetupRoutes() {

	s.engine.Post("/auth", handler.LoginUserHandler(s.userService))
	s.engine.Post("/users", handler.CreateUserHandler(s.userService))
	s.engine.Delete("/{userID}", middleware.ValidateJWTMiddleware(handler.DeleteUserHandler(s.userService)))
}
