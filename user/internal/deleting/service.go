package deleting

import (
	"context"
	"log"
	"net/http"

	"github.com/mata649/cqrs_on_aws/kit/response"
	user "github.com/mata649/cqrs_on_aws/user/internal"
)

type DeleteUserService struct {
	repository user.UserRepository
}

func NewDeleteUserService(repository user.UserRepository) DeleteUserService {
	return DeleteUserService{repository: repository}
}
func (u DeleteUserService) Delete(ctx context.Context, userID, currentUserID string) response.Response {
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
