package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mata649/cqrs_on_aws/auth"
	"github.com/mata649/cqrs_on_aws/response"
	"github.com/mata649/cqrs_on_aws/user/application"
	"github.com/mata649/cqrs_on_aws/user/domain"
)

var (
	headers = map[string]string{
		"Content-Type": "application/json",
	}
	badRequestResponse = events.APIGatewayProxyResponse{
		StatusCode: 409,
		Body:       `{"message" : "Bad Request"}`,
		Headers:    headers,
	}

	internalServerErrorResponse = events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `"message":"Internal server error"`,
		Headers:    headers,
	}
)

func CreateUserHandler(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	request := domain.CreateUserRequest{}
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Println(err)
		return badRequestResponse
	}
	resp := application.Create(ctx, request)
	return response.HandleResponse(resp)
}

func LoginUserHandler(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	request := domain.LoginUserRequest{}
	err := json.Unmarshal([]byte(req.Body), &request)

	if err != nil {
		log.Println(err)
		return badRequestResponse
	}
	resp := application.Login(ctx, request)

	if resp.GetType() != http.StatusOK {
		return response.HandleResponse(resp)

	}
	userResp, ok := resp.GetValue().(domain.UserResponse)
	if !ok {
		log.Println("Error casting UserResponse")
		return internalServerErrorResponse
	}
	jwt, err := auth.GenerateJWT(userResp.ID)
	if err != nil {
		log.Println("Error generating JWT", err)
		return internalServerErrorResponse

	}
	resMap := map[string]interface{}{
		"id":       userResp.ID,
		"username": userResp.Username,
		"createAt": userResp.CreatedAt,
		"jwt":      jwt,
	}
	jsonStr, err := json.Marshal(resMap)
	if err != nil {
		log.Println(err)
		return internalServerErrorResponse
	}
	return events.APIGatewayProxyResponse{
		Body:       string(jsonStr),
		StatusCode: http.StatusOK,
		Headers:    headers,
	}

}

func DeleteUserHandler(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	request := domain.DeleteUserRequest{
		UserID:        req.PathParameters["userID"],
		CurrentUserID: req.Headers["CurrentUserID"],
	}
	resp := application.Delete(ctx, request)

	return response.HandleResponse(resp)
}

func ChangePasswordHandler(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {

	var request domain.ChangePasswordRequest
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Println(err)
		return badRequestResponse
	}
	request.CurrentUserID = req.Headers["CurrentUserID"]
	resp := application.ChangePassword(ctx, request)
	return response.HandleResponse(resp)

}

func UpdateUserHandler(ctx context.Context, req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var request domain.UpdateUserRequest
	request.CurrentUserID = req.Headers["CurrentUserID"]
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		log.Println(err)
		return badRequestResponse
	}
	resp := application.Update(ctx, request)
	return response.HandleResponse(resp)
}
