package handler

import (
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/command"
	"github.com/mata649/cqrs_on_aws/user/internal/deleting"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/response"
)

type DeleteUserRequest struct {
	CurrentUserID string
	UserID        string
}

func DeleteUserHandler(commandBus command.Bus) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := DeleteUserRequest{
			UserID:        r.URL.Query().Get("userID"),
			CurrentUserID: r.Header.Get("CurrentUserID"),
		}

		resp := commandBus.Dispatch(r.Context(), deleting.NewDeleteUserCommand(request.UserID, request.CurrentUserID))

		response.WriteResponse(resp.GetType(), resp.GetValue(), w)

	})

}
