package bootstrap

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/mata649/cqrs_on_aws/platform/adapter"
	"github.com/mata649/cqrs_on_aws/platform/auth"
	server "github.com/mata649/cqrs_on_aws/platform/server/user"
	"github.com/mata649/cqrs_on_aws/platform/storage/dynamo"
	service "github.com/mata649/cqrs_on_aws/service/user"
)

type Config struct {
	UserTable string `envconfig:"USER_TABLE"`
	KeySecret string `evnconfig:"KEY_SECRET"`
}

func Run() error {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))
	auth.SetupAuth(config.KeySecret)
	userRepo := dynamo.NewUserDynamoRepository(config.UserTable, sess)
	userService := service.NewUserService(userRepo)

	server := server.NewServer(userService)
	server.SetupRoutes()
	lambdaAdapter := adapter.NewLambdaAdapter(server)

	lambda.StartWithOptions(lambdaAdapter.Handle, lambda.WithContext(context.Background()))
	return nil
}
