package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/domain/command"
	"github.com/mata649/cqrs_on_aws/platform/server/request"
	"github.com/mata649/cqrs_on_aws/platform/server/response"
	"github.com/mata649/cqrs_on_aws/service/task/creating"
)

type CreateTaskRequest struct {
	ID          string `json:"taskID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"userID"`
}

func CreateTaskHandler(commandBus command.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := CreateTaskRequest{}
		err := request.Binding(&req, r)
		if err != nil {
			log.Println(err)
			response.WriteBindingErrorResponse(err.Error(), w)
			return
		}
		resp := commandBus.Dispatch(r.Context(), creating.NewCreateTaskCommand(
			req.ID,
			req.Title,
			req.Description,
			req.UserID,
			time.Now(),
		))

		response.WriteResponse(resp.GetType(), resp.GetValue(), w)
	}

}
