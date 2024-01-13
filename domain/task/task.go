package task

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mata649/cqrs_on_aws/domain/event"
)

type TaskID struct {
	value string
}

func NewTaskID(value string) (TaskID, error) {
	_, err := uuid.Parse(value)
	if err != nil {
		return TaskID{}, errors.New("ID: is not valid")
	}
	return TaskID{value: value}, nil
}

func (t TaskID) String() string {
	return t.value
}

type TaskTitle struct {
	value string
}

func NewTaskTitle(value string) (TaskTitle, error) {
	if value == "" {
		return TaskTitle{}, errors.New("Title: is not valid")
	}
	return TaskTitle{value: value}, nil
}

func (t TaskTitle) String() string {
	return t.value
}

type TaskDescription struct {
	value string
}

func NewTaskDescription(value string) (TaskDescription, error) {
	if value == "" {
		return TaskDescription{}, errors.New("Description: is not valid")
	}
	return TaskDescription{value: value}, nil
}

func (t TaskDescription) String() string {
	return t.value
}

type TaskCreatedAt struct {
	value time.Time
}

func NewTaskCreatedAt(value time.Time) (TaskCreatedAt, error) {
	if value.IsZero() || value.After(time.Now()) {
		return TaskCreatedAt{}, errors.New("CreatedAt: is not valid")

	}

	return TaskCreatedAt{value: value}, nil
}

func (t TaskCreatedAt) Time() time.Time {
	return t.value
}

type TaskUserID struct {
	value string
}

func NewTaskUserID(value string) (TaskUserID, error) {
	_, err := uuid.Parse(value)
	if err != nil {
		return TaskUserID{}, errors.New("UserID: is not valid")
	}
	return TaskUserID{value: value}, nil
}
func (t TaskUserID) String() string {
	return t.value
}

type Task struct {
	id          TaskID
	title       TaskTitle
	description TaskDescription
	createdAt   TaskCreatedAt
	userID      TaskUserID
	events      []event.Event
}

func NewTask(id, title, description string, createdAt time.Time, userID string) (Task, error) {
	errs := []error{}
	idVO, err := NewTaskID(id)
	if err != nil {
		errs = append(errs, err)
	}
	titleVO, err := NewTaskTitle(title)
	if err != nil {
		errs = append(errs, err)
	}
	descriptionVO, err := NewTaskDescription(description)
	if err != nil {
		errs = append(errs, err)
	}
	createdAtVO, err := NewTaskCreatedAt(createdAt)
	if err != nil {
		errs = append(errs, err)
	}
	userIDVO, err := NewTaskUserID(userID)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return Task{}, errs[0]
	}
	task := Task{
		id:          idVO,
		title:       titleVO,
		description: descriptionVO,
		createdAt:   createdAtVO,
		userID:      userIDVO,
	}
	task.Record(NewTaskCreatedEvent(id, title, description, createdAt, userID))
	return task, nil
}

func (u *Task) ID() TaskID {
	return u.id
}

func (u *Task) Title() TaskTitle {
	return u.title
}

func (u *Task) Description() TaskDescription {
	return u.description
}

func (u *Task) CreatedAt() TaskCreatedAt {
	return u.createdAt
}

func (u *Task) UserID() TaskUserID {
	return u.userID
}

func (u *Task) Record(e event.Event) {
	u.events = append(u.events, e)
}

func (u *Task) PullEvents() []event.Event {
	events := u.events
	u.events = []event.Event{}
	return events
}

type TaskRepository interface {
	Close() error
	Create(ctx context.Context, task Task) error
	GetByID(ctx context.Context, id string) (Task, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, task Task) error
}
