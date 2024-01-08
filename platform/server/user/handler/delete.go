package handler

import (
	"net/http"

	"github.com/mata649/cqrs_on_aws/platform/server/response"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

type DeleteUserRequest struct {
	CurrentUserID string
	UserID        string
}

func DeleteUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := DeleteUserRequest{
			UserID:        r.URL.Query().Get("userID"),
			CurrentUserID: r.Header.Get("CurrentUserID"),
		}

		resp := userService.Delete(r.Context(), request.CurrentUserID, request.UserID)

		response.WriteResponse(resp.GetType(), resp.GetValue(), w)

	}
}
