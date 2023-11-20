package domain

import (
	"context"
	"log"
	"net/http"

	"github.com/mata649/cqrs_on_aws/response"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (service UserService) Create(ctx context.Context, request CreateUserRequest) response.Response {
	user, err := request.Validate()
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponseFailure(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := service.userRepository.GetByUsername(ctx, user.Username)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	if userFound.ID != "" {
		return response.NewResponseFailure(http.StatusNotFound, "Username already exists")
	}
	passwordHashed, err := hashPassword(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	user.Password = passwordHashed
	err = service.userRepository.Create(ctx, user)
	if err != nil {
		log.Println("Error creating user", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}

	return response.NewResponseSuccessful(http.StatusCreated, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

func (service UserService) Login(ctx context.Context, request LoginUserRequest) response.Response {
	user, err := request.Validate()
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponseFailure(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := service.userRepository.GetByUsername(ctx, user.Username)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	if userFound.ID == "" {
		return response.NewResponseFailure(http.StatusUnauthorized, "Invalid credentials")
	}
	if checkPasswordHash(user.Password, userFound.Password) == false {
		return response.NewResponseFailure(http.StatusUnauthorized, "Invalid credentials")
	}
	return response.NewResponseSuccessful(http.StatusOK, UserResponse{
		ID:        userFound.ID,
		Username:  userFound.Username,
		CreatedAt: userFound.CreatedAt,
	})
}

func (service UserService) Delete(ctx context.Context, request DeleteUserRequest) response.Response {
	err := request.Validate()
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponseFailure(400, "Bad Request")
	}
	userFound, err := service.userRepository.GetByID(ctx, request.UserID)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}

	if userFound.ID == "" {
		return response.NewResponseFailure(http.StatusNotFound, "The user doesn't exist")
	}
	if userFound.ID != request.CurrentUserID {
		return response.NewResponseFailure(http.StatusUnauthorized, "Unauthorized")
	}
	err = service.userRepository.Delete(ctx, userFound.ID)
	if err != nil {
		log.Println("Error deletin user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	return response.NewResponseSuccessful(http.StatusOK, UserResponse{
		ID:        userFound.ID,
		Username:  userFound.Username,
		CreatedAt: userFound.CreatedAt,
	})
}
func (s UserService) ChangePassword(ctx context.Context, request ChangePasswordRequest) response.Response {
	err := request.Validate()
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponseFailure(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := s.userRepository.GetByID(ctx, request.CurrentUserID)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	if userFound.ID == "" {
		return response.NewResponseFailure(http.StatusNotFound, "The user doesn't exist")
	}
	if !checkPasswordHash(request.OldPassword, userFound.Password) {
		return response.NewResponseFailure(http.StatusUnauthorized, "The old password is not correct")
	}
	passwordHashed, err := hashPassword(request.NewPassword)
	if err != nil {
		log.Println("Error hashing password:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	userFound.Password = passwordHashed
	err = s.userRepository.Update(ctx, userFound)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}

	return response.NewResponseSuccessful(http.StatusOK, UserResponse{
		ID:        userFound.ID,
		Username:  userFound.Username,
		CreatedAt: userFound.CreatedAt,
	})

}

func (s UserService) Update(ctx context.Context, request UpdateUserRequest) response.Response {
	user, err := request.Validate()
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponseFailure(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := s.userRepository.GetByID(ctx, user.ID)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}
	if userFound.ID == "" {
		return response.NewResponseFailure(http.StatusNotFound, "The user doesn't exist")
	}
	if userFound.ID != request.CurrentUserID {
		return response.NewResponseFailure(http.StatusUnauthorized, "Unauthorized")
	}
	userFound.Username = user.Username
	err = s.userRepository.Update(ctx, userFound)
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponseFailure(http.StatusInternalServerError, "Internal Server Error")
	}

	return response.NewResponseSuccessful(http.StatusOK, UserResponse{
		ID:        userFound.ID,
		Username:  userFound.Username,
		CreatedAt: userFound.CreatedAt,
	})

}
