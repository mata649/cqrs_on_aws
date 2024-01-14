package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mata649/cqrs_on_aws/domain/event"
)

type EventBus struct {
	sns      *sns.SNS
	exchange string
	sqs      *sqs.SQS
}

func NewEventBus(sess *session.Session, exchange string) (EventBus, error) {
	sns := sns.New(sess)
	availableTopics, _ := sns.ListTopics(nil)

	var requiredTopicArn string

	for _, t := range availableTopics.Topics {
		if strings.Contains(*t.TopicArn, exchange) {
			requiredTopicArn = exchange
		}
	}
	if requiredTopicArn == "" {
		return EventBus{}, fmt.Errorf("topic %s was not found", exchange)

	}

	return EventBus{
		sns:      sns,
		sqs:      sqs.New(sess),
		exchange: requiredTopicArn,
	}, nil
}

func (e EventBus) Publish(ctx context.Context, events []event.Event) error {

	for _, event := range events {
		jsonBody, err := json.Marshal(event)
		if err != nil {
			return err
		}
		_, err = e.sns.Publish(&sns.PublishInput{
			Message:  aws.String(string(jsonBody)),
			TopicArn: aws.String(string(e.exchange)),
			MessageAttributes: map[string]*sns.MessageAttributeValue{
				"event_type": {
					DataType:    aws.String("String"),
					StringValue: aws.String(string(event.Type())),
				},
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (e EventBus) Subscribe(eventType event.Type, handler event.Handler) {

}
