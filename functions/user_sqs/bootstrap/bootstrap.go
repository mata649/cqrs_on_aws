package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kelseyhightower/envconfig"
	"github.com/mata649/cqrs_on_aws/platform/storage/dynamo"
	"github.com/mata649/cqrs_on_aws/service/user/increasing"
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
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))
	userRepo := dynamo.NewUserDynamoRepository(config.UserTable, sess)

	_ = increasing.NewTaskCounterIncreaserService(userRepo)

	return nil
}
