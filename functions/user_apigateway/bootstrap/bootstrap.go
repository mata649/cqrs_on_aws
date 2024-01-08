package bootstrap

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/mata649/cqrs_on_aws/platform/adapter"
	"github.com/mata649/cqrs_on_aws/platform/auth"
	server "github.com/mata649/cqrs_on_aws/platform/server/user"
	dynamo "github.com/mata649/cqrs_on_aws/platform/storage/dynamo/user"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

type Config struct {
	UserTable          string `envconfig:"USER_TABLE"`
	GoogleClientSecret string `envconfig:"GOOGLE_KEY"`
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
	KeySecret          string `envconfig:"KEY_SECRET"`
}

func Run() error {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return err
	}
	auth.SetupAuth(config.GoogleClientID, config.GoogleClientSecret, config.KeySecret)
	userRepo := dynamo.NewUserDynamoRepository(config.UserTable)
	userService := service.NewUserService(userRepo)

	server := server.NewServer(userService)
	server.SetupRoutes()
	lambdaAdapter := adapter.NewLambdaAdapter(server)

	lambda.StartWithOptions(lambdaAdapter.Handle, lambda.WithContext(context.Background()))
	return nil
}
