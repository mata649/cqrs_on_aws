package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/kit/command"
	"github.com/mata649/cqrs_on_aws/user/internal/creating"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/request"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/response"
)

type CreateUserRequest struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	CreateAt time.Time `json:"createAt"`
}

func CreateUserHandler(commandBus command.Bus) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := CreateUserRequest{}
		err := request.Binding(&req, r)
		if err != nil {
			log.Println(err)
			response.WriteResponse(http.StatusBadRequest, "Bad Request", w)
			return
		}
		resp := commandBus.Dispatch(r.Context(), creating.NewCreateUserCommand(req.ID, req.Username, req.Password, req.CreateAt))

		response.WriteResponse(resp.GetType(), resp.GetValue(), w)

	})
}
