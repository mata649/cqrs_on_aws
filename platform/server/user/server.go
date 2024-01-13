package server

import (
	"log"
	"net/http"

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
	s.engine.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL)
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	})
	s.engine.Post("/auth", handler.LoginUserHandler(s.userService))
	s.engine.Post("/users", handler.CreateUserHandler(s.userService))
	s.engine.Get("/users", handler.GetUserByIDHandler(s.userService))
	s.engine.Delete("/users/{userID}", middleware.ValidateJWTMiddleware(handler.DeleteUserHandler(s.userService)))
}
