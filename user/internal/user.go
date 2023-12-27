package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mata649/cqrs_on_aws/kit/event"
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

type UserUsername struct {
	value string
}

func NewUserUsername(value string) (UserUsername, error) {
	if len(value) < 3 || len(value) > 40 {
		return UserUsername{}, errors.New("Username: is not valid")
	}
	return UserUsername{value: value}, nil

}
func (u UserUsername) String() string {
	return u.value
}

type UserPassword struct {
	value string
}

func NewUserPassword(value string) (UserPassword, error) {
	if len(value) < 8 || len(value) > 64 {
		return UserPassword{}, errors.New("Password: is not valid")
	}
	return UserPassword{value: value}, nil
}
func (u UserPassword) String() string {
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

type User struct {
	id        UserID
	username  UserUsername
	password  UserPassword
	createdAt UserCreatedAt
	events    []event.Event
}

func NewUser(id, username, password string, createdAt time.Time) (User, error) {
	errs := []error{}
	idVO, err := NewUserID(id)
	if err != nil {
		errs = append(errs, err)
	}
	usernameVO, err := NewUserUsername(username)
	if err != nil {
		errs = append(errs, err)
	}
	passwordVO, err := NewUserPassword(password)
	if err != nil {
		errs = append(errs, err)
	}
	createdAtVO, err := NewUserCreatedAt(createdAt)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return User{}, errors.Join(errs...)
	}
	return User{
		id:        idVO,
		username:  usernameVO,
		password:  passwordVO,
		createdAt: createdAtVO,
	}, nil
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Username() UserUsername {
	return u.username
}
func (u User) Password() UserPassword {
	return u.password
}
func (u User) CreatedAt() UserCreatedAt {
	return u.createdAt
}
func (u *User) SetHashedPassword(password string) error {
	userPassword, err := NewUserPassword(password)
	if err != nil {
		return err
	}
	u.password = userPassword
	return nil
}
func (u *User) UpdateUsername(username UserUsername) {

}
func (u *User) PullEvents() []event.Event {
	events := u.events
	u.events = []event.Event{}
	return events
}

func (u User) Record(event event.Event) {
	u.events = append(u.events, event)
}

type UserRepository interface {
	Close()
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	Get(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user User) error
}
