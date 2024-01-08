package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/domain/response"
	"github.com/mata649/cqrs_on_aws/domain/user"
)

type UserService struct {
	repository user.UserRepository
}

func NewUserService(repository user.UserRepository) UserService {
	return UserService{
		repository: repository,
	}
}

type UserResponse struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

func (u UserService) Create(ctx context.Context, id, email string, createdAt time.Time) response.Response {
	user, err := user.NewUser(id, email, createdAt, 0)

	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := u.repository.GetByEmail(ctx, user.Email().String())
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}

	if userFound.ID().String() != "" {
		return response.NewResponse(http.StatusNotFound, "Email already exists")
	}

	err = u.repository.Create(ctx, user)

	if err != nil {
		log.Println("Error creating user", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}

	return response.NewResponse(http.StatusCreated, "User created successfully")
}

func (u UserService) Delete(ctx context.Context, userID, currentUserID string) response.Response {
	userIDVO, err := user.NewUserID(userID)
	if err != nil {
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	currentUserIDVO, err := user.NewUserID(currentUserID)

	if err != nil {
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}

	userFound, err := u.repository.GetByID(ctx, userIDVO.String())
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}

	if userFound.ID().String() == "" {
		return response.NewResponse(http.StatusNotFound, "The user doesn't exist")
	}
	if userFound.ID().String() != currentUserIDVO.String() {
		return response.NewResponse(http.StatusUnauthorized, "Unauthorized")
	}
	err = u.repository.Delete(ctx, userFound.ID().String())
	if err != nil {
		log.Println("Error deletin user:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	return response.NewResponse(http.StatusOK, "User deleted successfully")
}

func (u UserService) Authenticate(ctx context.Context, email string) response.Response {
	usernameVO, err := user.NewUserEmail(email)
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}

	userFound, err := u.repository.GetByEmail(ctx, usernameVO.String())
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	if userFound.ID().String() == "" {
		return response.NewResponse(http.StatusUnauthorized, "Invalid credentials")
	}

	return response.NewResponse(http.StatusOK, UserResponse{
		ID:        userFound.ID().String(),
		Email:     userFound.Email().String(),
		CreatedAt: userFound.CreatedAt().Time(),
	})
}
