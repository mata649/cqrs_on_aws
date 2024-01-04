package middleware

import (
	"log"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/platform/server/auth"
	"github.com/mata649/cqrs_on_aws/kit/platform/server/response"
)

func ValidateJWTMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Header["Authorization"]
		if !ok {
			response.WriteResponse(http.StatusUnauthorized, "Unauthorized", w)
			return
		}

		claims, err := auth.GetClaimsFromToken(token[0])
		if err != nil {
			log.Println("Error getting claims:", err)
			response.WriteResponse(http.StatusUnauthorized, "Unauthorized", w)
		}
		r.Header["CurrentUserID"] = []string{claims.ID}

		next.ServeHTTP(w, r)
	})
}
