package task

import (
	"time"

	"github.com/mata649/cqrs_on_aws/domain/event"
)

const TaskCreatedEventEventType event.Type = "events.task.created"

type TaskCreatedEvent struct {
	event.BaseEvent
	id          string
	title       string
	description string
	createdAt   time.Time
	userID      string
}

func NewTaskCreatedEvent(id, title, description string, createdAt time.Time, userID string) TaskCreatedEvent {
	return TaskCreatedEvent{
		id:          id,
		title:       title,
		description: description,
		createdAt:   createdAt,
		userID:      userID,
	}
}

func (e TaskCreatedEvent) Type() event.Type {
	return TaskCreatedEventEventType
}
func (e TaskCreatedEvent) ID() string {
	return e.id
}

func (e TaskCreatedEvent) Title() string {
	return e.title
}

func (e TaskCreatedEvent) Description() string {
	return e.description
}

func (e TaskCreatedEvent) CreatedAt() time.Time {
	return e.createdAt
}

func (e TaskCreatedEvent) UserID() string {
	return e.userID
}
