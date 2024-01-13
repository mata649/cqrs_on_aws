package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mata649/cqrs_on_aws/platform/server/auth"
	"github.com/mata649/cqrs_on_aws/platform/server/request"
	"github.com/mata649/cqrs_on_aws/platform/server/response"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

type LoginUserRequest struct {
	Email string `json:"email"`
}

func LoginUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := LoginUserRequest{}
		err := request.Binding(&req, r)
		if err != nil {
			log.Println(err)
			response.WriteBindingErrorResponse(err.Error(), w)
			return
		}
		resp := userService.Authenticate(r.Context(), req.Email)

		if resp.GetType() != http.StatusOK {
			response.WriteResponse(resp.GetType(), resp.GetValue(), w)
			return

		}
		userResp, ok := resp.GetValue().(service.UserResponse)
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
			"email":    userResp.Email,
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
