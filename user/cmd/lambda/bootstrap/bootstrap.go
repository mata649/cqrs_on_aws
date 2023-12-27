package bootstrap

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/mata649/cqrs_on_aws/user/internal/authenticating"
	"github.com/mata649/cqrs_on_aws/user/internal/creating"
	"github.com/mata649/cqrs_on_aws/user/internal/deleting"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/adapter"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/bus/inmemory"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/storage/dynamo"
)

type Config struct {
	UserTable string `envconfig:"USER_TABLE"`
}

func Run() error {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return err
	}
	userRepo := dynamo.NewUserDynamoRepository(config.UserTable)
	deleteUserService := deleting.NewDeleteUserService(userRepo)
	creatingService := creating.NewCreateUserService(userRepo)
	authenticatingService := authenticating.NewAuthenticateUserService(userRepo)

	commandBus := inmemory.NewCommandBus()
	commandBus.Register(creating.CreateUserCommandType, creating.NewCreateUserCommandHandler(creatingService))
	commandBus.Register(deleting.DeleteUserCommandType, deleting.NewDeleteUserCommandHandler(deleteUserService))

	queryBus := inmemory.NewQueryBus()
	queryBus.Register(authenticating.AuthenticateUserQueryType, authenticating.NewAuthenticateUserQueryHandler(authenticatingService))

	server := server.NewServer(commandBus, queryBus)
	lambdaAdapter := adapter.NewLambdaAdapter(server)

	lambda.StartWithOptions(lambdaAdapter.Handle, lambda.WithContext(context.Background()))
	return nil
}
