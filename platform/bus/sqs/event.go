package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/mata649/cqrs_on_aws/domain/event"
)

type EventBus struct {
	svc *sns.SNS
}

func NewEventBus(sess *session.Session) EventBus {
	return EventBus{
		sns.New(sess),
	}
}

func (e EventBus) Publish(ctx context.Context, event event.Event) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = e.svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(jsonBody)),
		TopicArn: aws.String(string(event.Type())),
	})
	if err != nil {
		return err
	}
	return nil
}

func (e EventBus) Subscribe(eventType event.Type, handler event.Handler) {

}
