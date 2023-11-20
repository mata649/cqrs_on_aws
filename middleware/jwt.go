package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mata649/cqrs_on_aws/auth"
)

type HandlerFunc func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

var unauthorizedResponse = events.APIGatewayProxyResponse{
	StatusCode: http.StatusUnauthorized,
	Body:       `{"message":"Unauthorized"}`,
	Headers: map[string]string{
		"Content-Type": "application/json",
	},
}

func ValidateJWTMiddleware(h func(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) func(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

		token, ok := req.Headers["Authorization"]

		if !ok {
			return unauthorizedResponse
		}
		claims, err := auth.GetClaimsFromToken(token)
		if err != nil {
			log.Println("Error getting claims:", err)
			return unauthorizedResponse
		}
		req.Headers["CurrentUserID"] = claims.ID
		return h(ctx, req)
	}
}
