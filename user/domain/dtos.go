package domain

import (
	"log"
	"time"

	"github.com/segmentio/ksuid"
	"gopkg.in/validator.v2"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"min=3,max=40"`
	Password string `json:"password" validate:"min=8,max=64"`
}

func (r CreateUserRequest) Validate() (User, error) {
	err := validator.Validate(r)
	if err != nil {
		return User{}, err
	}
	user := User{}
	user.CreatedAt = time.Now().UTC()

	uuid, err := ksuid.NewRandom()
	if err != nil {
		log.Println(err)
		return user, err
	}
	user.ID = uuid.String()
	user.Password = r.Password
	user.Username = r.Username
	return user, nil
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

func (r LoginUserRequest) Validate() (User, error) {
	err := validator.Validate(r)
	if err != nil {
		return User{}, err
	}
	return User{
		Username: r.Username,
		Password: r.Password,
	}, nil
}

type DeleteUserRequest struct {
	CurrentUserID string
	UserID        string
}

func (r DeleteUserRequest) Validate() error {

	_, err := ksuid.Parse(r.CurrentUserID)
	if err != nil {
		log.Println("CurrentUserID: is not valid")
		return err
	}
	_, err = ksuid.Parse(r.UserID)
	if err != nil {
		log.Println("UserID: is not valid")
		return err
	}
	return nil

}

type ChangePasswordRequest struct {
	CurrentUserID string
	OldPassword   string `json:"oldPassword" validate:"nonzero"`
	NewPassword   string `json:"newPassword" validate:"min=8,max=64"`
}

func (r ChangePasswordRequest) Validate() error {
	err := validator.Validate(r)
	if err != nil {
		return err
	}
	_, err = ksuid.Parse(r.CurrentUserID)
	if err != nil {
		log.Println("CurrentUserID: is not valid")
		return err
	}
	return nil
}

type UpdateUserRequest struct {
	CurrentUserID string
	UserID        string
	Username      string `json:"username" validate:"min=3,max=40"`
}

func (r UpdateUserRequest) Validate() (User, error) {
	err := validator.Validate(r)
	if err != nil {
		return User{}, err
	}
	_, err = ksuid.Parse(r.CurrentUserID)
	if err != nil {
		log.Println("CurrentUserID: is not valid")
		return User{}, err
	}
	return User{
		ID:       r.UserID,
		Username: r.Username,
	}, nil
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}
