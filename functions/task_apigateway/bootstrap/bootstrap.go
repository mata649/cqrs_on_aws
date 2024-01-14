package bootstrap

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/mata649/cqrs_on_aws/platform/adapter"
	"github.com/mata649/cqrs_on_aws/platform/auth"
	"github.com/mata649/cqrs_on_aws/platform/bus/aws"
	"github.com/mata649/cqrs_on_aws/platform/bus/inmemory"
	server "github.com/mata649/cqrs_on_aws/platform/server/task"
	"github.com/mata649/cqrs_on_aws/platform/storage/dynamo"
	"github.com/mata649/cqrs_on_aws/service/task/creating"
)

type Config struct {
	UserTable string `envconfig:"TASK_TABLE"`
	KeySecret string `evnconfig:"KEY_SECRET"`
	Exchange  string `envconfig:"EXCHANGE"`
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
	eventBus, err := aws.NewEventBus(sess, config.Exchange)
	if err != nil {
		return err
	}
	auth.SetupAuth(config.KeySecret)
	taskRepo := dynamo.NewTaskDynamoRepository(config.UserTable, sess)

	taskService := creating.NewCreateTaskService(taskRepo, eventBus)

	commandBus := inmemory.NewCommandBus()
	commandBus.Register(creating.CreateTaskCommandType, creating.NewCreateTaskCommandHandler(taskService))

	server := server.NewServer(commandBus)
	server.SetupRoutes()

	lambdaAdapter := adapter.NewLambdaAdapter(server)
	lambda.StartWithOptions(lambdaAdapter.Handle, lambda.WithContext(context.Background()))
	return nil
}
