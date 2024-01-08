package user

import "github.com/mata649/cqrs_on_aws/domain/event"

const TaskCreatedEventEventType event.Type = "events.task.created"

type TaskCreatedEvent struct {
	event.BaseEvent
	userID string
}

func (e TaskCreatedEvent) UserID() string {
	return e.userID

}

func (e TaskCreatedEvent) Type() event.Type {
	return TaskCreatedEventEventType

}
