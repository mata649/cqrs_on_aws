package handler

import (
	"net/http"

	"github.com/mata649/cqrs_on_aws/user/internal/platform/server/response"
)

func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.WriteResponse(http.StatusOK, "OK", w)
	}
}
