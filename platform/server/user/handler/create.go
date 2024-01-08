package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/platform/server/request"
	"github.com/mata649/cqrs_on_aws/platform/server/response"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

type CreateUserRequest struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
}

func CreateUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := CreateUserRequest{}
		err := request.Binding(&req, r)
		if err != nil {
			log.Println(err)
			response.WriteBindingErrorResponse(err.Error(), w)
			return
		}
		resp := userService.Create(r.Context(), req.UserID, req.Email, time.Now())

		response.WriteResponse(resp.GetType(), resp.GetValue(), w)
	}
}
