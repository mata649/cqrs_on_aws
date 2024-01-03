package creating

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/kit/response"
	user "github.com/mata649/cqrs_on_aws/user/internal"
	"github.com/mata649/cqrs_on_aws/user/internal/utils"
)

type CreateUserService struct {
	repository user.UserRepository
}

func NewCreateUserService(repository user.UserRepository) CreateUserService {
	return CreateUserService{repository: repository}
}

func (u CreateUserService) Create(ctx context.Context, id, username, password string, createdAt time.Time) response.Response {
	user, err := user.NewUser(id, username, password, createdAt)
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := u.repository.GetByUsername(ctx, user.Username().String())
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}

	if userFound.ID().String() != "" {
		return response.NewResponse(http.StatusNotFound, "Username already exists")
	}
	passwordHashed, err := utils.HashPassword(user.Password().String())
	if err != nil {
		log.Println("Error hashing password:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	if user.SetHashedPassword(passwordHashed); err != nil {
		log.Println("Error setting hashed password:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}

	err = u.repository.Create(ctx, user)

	if err != nil {
		log.Println("Error creating user", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}

	return response.NewResponse(http.StatusCreated, "User created successfully")
}
