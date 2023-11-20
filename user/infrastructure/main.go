package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"

	"github.com/mata649/cqrs_on_aws/middleware"
	"github.com/mata649/cqrs_on_aws/user/database"
	"github.com/mata649/cqrs_on_aws/user/repository"
)

type Config struct {
	UserTable string `envconfig:"USER_TABLE"`
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case http.MethodPost:
		if req.Path == "/auth" {
			return LoginUserHandler(ctx, req), nil
		}
		return CreateUserHandler(ctx, req), nil
	case http.MethodPut:
		if req.Path == "/changePassword" {
			return middleware.ValidateJWTMiddleware(ChangePasswordHandler)(ctx, req), nil
		}
		return middleware.ValidateJWTMiddleware(UpdateUserHandler)(ctx, req), nil
	case http.MethodDelete:
		return middleware.ValidateJWTMiddleware(DeleteUserHandler)(ctx, req), nil
	default:
		log.Println("Invalid HTTPMethod")
		return events.APIGatewayProxyResponse{
			Body:       "Method not supported",
			StatusCode: 405,
		}, nil
	}

}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetUserRepository(database.NewUserDynamoRepository(config.UserTable))

	lambda.Start(router)

}
