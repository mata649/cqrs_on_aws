package authenticating

import (
	"context"
	"log"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/response"
	user "github.com/mata649/cqrs_on_aws/user/internal"
	"github.com/mata649/cqrs_on_aws/user/internal/utils"
)

type AuthenticateUserService struct {
	repository user.UserRepository
}

func NewAuthenticateUserService(repository user.UserRepository) AuthenticateUserService {
	return AuthenticateUserService{repository: repository}

}

func (u AuthenticateUserService) Authenticate(ctx context.Context, username, password string) response.Response {
	usernameVO, err := user.NewUserUsername(username)
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	passwordVO, err := user.NewUserPassword(password)
	if err != nil {
		log.Println("Request error:", err)
		return response.NewResponse(http.StatusBadRequest, response.ParseErrorResponse(err.Error()))
	}
	userFound, err := u.repository.GetByUsername(ctx, usernameVO.String())
	if err != nil {
		log.Println("Error getting user:", err)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	if userFound.ID().String() == "" {
		return response.NewResponse(http.StatusUnauthorized, "Invalid credentials")
	}
	if utils.CheckPasswordHash(passwordVO.String(), userFound.Password().String()) == false {
		return response.NewResponse(http.StatusUnauthorized, "Invalid credentials")
	}
	return response.NewResponse(http.StatusOK, AuthenticateUserQueryResponse{
		ID:        userFound.ID().String(),
		Username:  userFound.Username().String(),
		CreatedAt: userFound.CreatedAt().Time(),
	})
}
