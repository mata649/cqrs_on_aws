package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mata649/cqrs_on_aws/domain/event"
)

type EventBus struct {
	sns             *sns.SNS
	availableTopics []string
	sqs             *sqs.SQS
}

func NewEventBus(sess *session.Session) EventBus {
	sns := sns.New(sess)
	availableTopics, _ := sns.ListTopics(nil)
	availableTopicsArns := make([]string, len(availableTopics.Topics))
	for i, t := range availableTopics.Topics {
		availableTopicsArns[i] = *t.TopicArn
	}

	return EventBus{
		sns:             sns,
		sqs:             sqs.New(sess),
		availableTopics: availableTopicsArns,
	}
}

func (e EventBus) Publish(ctx context.Context, event event.Event) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return err
	}

	var requiredTopicArn string

	for _, t := range e.availableTopics {
		if strings.Contains(t, string(event.Type())) {
			requiredTopicArn = t
		}
	}
	if len(requiredTopicArn) > 0 {
		_, err = e.sns.Publish(&sns.PublishInput{
			Message:  aws.String(string(jsonBody)),
			TopicArn: aws.String(string(event.Type())),
		})

		if err != nil {
			return err
		}
		return nil
	}

	return errors.New(fmt.Sprintf("Topic %s was not found \n", string(event.Type())))
}

func (e EventBus) Subscribe(eventType event.Type, handler event.Handler) {
	e.sqs.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(string(eventType)),
	})
}
