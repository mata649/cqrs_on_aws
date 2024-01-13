package handler

import (
	"net/http"

	"github.com/mata649/cqrs_on_aws/platform/server/response"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

func GetUserByIDHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userID")
		resp := userService.GetByID(r.Context(), userID)
		response.WriteResponse(resp.GetType(), resp.GetValue(), w)
	}
}
