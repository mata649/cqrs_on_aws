package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Bus interface {
	Publish(context.Context, []Event) error
	Subscribe(Type, Handler)
}
type Handler interface {
	Handle(context.Context, Event) error
}
type Type string

type Event interface {
	ID() string
	AggregateID() string
	OccurredOn() time.Time
	Type() Type
}

type BaseEvent struct {
	eventID     string
	aggregateID string
	occurredOn  time.Time
}

func (b BaseEvent) ID() string {
	return b.eventID

}
func (b BaseEvent) AggregateID() string {
	return b.aggregateID
}
func (b BaseEvent) OccurredOn() time.Time {
	return b.occurredOn
}

func NewBaseEvent(aggregateID string) BaseEvent {
	return BaseEvent{
		aggregateID: aggregateID,
		eventID:     uuid.New().String(),
		occurredOn:  time.Now(),
	}

}
