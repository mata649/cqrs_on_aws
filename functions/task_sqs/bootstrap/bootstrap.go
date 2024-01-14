package bootstrap

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/mata649/cqrs_on_aws/platform/storage/dynamo"
	"github.com/mata649/cqrs_on_aws/service/user/increasing"
)

type Config struct {
	UserTable string `envconfig:"USER_TABLE"`
}

var increaserService increasing.TaskCounterIncreaserService

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		err := increaserService.IncreaseTasksCreatedCounter(ctx, record.Body)
		if err != nil {
			return err
		}
	}
	return nil
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
	userRepo := dynamo.NewUserDynamoRepository(config.UserTable, sess)

	_ = increasing.NewTaskCounterIncreaserService(userRepo)

	lambda.StartWithOptions(Handler, lambda.WithContext(context.Background()))

	return nil
}
