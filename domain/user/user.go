package user

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/mata649/cqrs_on_aws/domain/event"
)

type UserID struct {
	value string
}

func NewUserID(value string) (UserID, error) {
	_, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, errors.New("ID: is not valid")
	}
	return UserID{value: value}, nil
}
func (u UserID) String() string {
	return u.value
}

type UserEmail struct {
	value string
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserEmail(value string) (UserEmail, error) {
	if !isEmailValid(value) {
		return UserEmail{}, errors.New("Email: is not valid")
	}
	return UserEmail{value: value}, nil

}
func (u UserEmail) String() string {
	return u.value
}

type UserCreatedAt struct {
	value time.Time
}

func (u UserCreatedAt) Time() time.Time {
	return u.value
}
func NewUserCreatedAt(value time.Time) (UserCreatedAt, error) {
	if value.IsZero() || value.After(time.Now()) {
		return UserCreatedAt{}, errors.New("CreatedAt: is not valid")

	}

	return UserCreatedAt{value: value}, nil
}

type UserTasksCreated struct {
	value uint
}

func (u UserTasksCreated) Uint() uint {
	return u.value
}
func NewUserTasksCreated(value uint) (UserTasksCreated, error) {
	if value < 0 {
		return UserTasksCreated{}, errors.New("TasksCreated: is not valid")

	}

	return UserTasksCreated{value: value}, nil
}

type User struct {
	id           UserID
	email        UserEmail
	createdAt    UserCreatedAt
	events       []event.Event
	tasksCreated UserTasksCreated
}

func NewUser(id, email string, createdAt time.Time, tasksCreated uint) (User, error) {
	errs := []error{}
	idVO, err := NewUserID(id)
	if err != nil {
		errs = append(errs, err)
	}
	emailVO, err := NewUserEmail(email)
	if err != nil {
		errs = append(errs, err)
	}

	createdAtVO, err := NewUserCreatedAt(createdAt)
	if err != nil {
		errs = append(errs, err)
	}
	tasksCreatedVO, err := NewUserTasksCreated(tasksCreated)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return User{}, errors.Join(errs...)
	}

	user := User{
		id:           idVO,
		email:        emailVO,
		createdAt:    createdAtVO,
		tasksCreated: tasksCreatedVO,
	}

	return user, nil
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Email() UserEmail {
	return u.email
}

func (u User) CreatedAt() UserCreatedAt {
	return u.createdAt
}
func (u User) TasksCreated() UserTasksCreated {
	return u.tasksCreated
}
func (u User) IncreaseTasksCreatedCounter() {
	u.tasksCreated.value++
}

type UserRepository interface {
	Close()
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user User) error
}
