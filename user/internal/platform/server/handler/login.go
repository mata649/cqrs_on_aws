package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/query"
	"github.com/mata649/cqrs_on_aws/user/internal/authenticating"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/auth"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/request"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/response"
)

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginUserHandler(queryBus query.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := LoginUserRequest{}
		err := request.Binding(&req, r)
		if err != nil {
			log.Println(err)
			response.WriteBindingErrorResponse(err.Error(), w)
			return
		}
		query := authenticating.NewAuthenticateUserQuery(req.Username, req.Password)

		resp := queryBus.Dispatch(r.Context(), query)

		if resp.GetType() != http.StatusOK {
			response.WriteResponse(resp.GetType(), resp.GetValue(), w)
			return

		}
		userResp, ok := resp.GetValue().(authenticating.AuthenticateUserQueryResponse)
		if !ok {
			log.Println("Error casting UserResponse")
			response.WriteResponse(http.StatusInternalServerError, "Internal Server Error", w)
			return
		}
		jwt, err := auth.GenerateJWT(userResp.ID)
		if err != nil {
			log.Println("Error generating JWT", err)
			response.WriteResponse(http.StatusInternalServerError, "Internal Server Error", w)
			return
		}
		resMap := map[string]interface{}{
			"id":       userResp.ID,
			"username": userResp.Username,
			"createAt": userResp.CreatedAt,
			"jwt":      jwt,
		}
		jsonStr, err := json.Marshal(resMap)
		if err != nil {
			log.Println(err)
			response.WriteResponse(http.StatusInternalServerError, "Internal Server Error", w)
			return
		}
		response.WriteResponse(http.StatusOK, string(jsonStr), w)

	}

}
