package authenticating

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mata649/cqrs_on_aws/kit/query"
	"github.com/mata649/cqrs_on_aws/kit/response"
)

const AuthenticateUserQueryType query.Type = "query.autenticating.user"

type AuthenticateUserQuery struct {
	Username string
	Password string
}

type AuthenticateUserQueryResponse struct {
	ID        string
	Username  string
	CreatedAt time.Time
}

func NewAuthenticateUserQuery(username, password string) AuthenticateUserQuery {
	return AuthenticateUserQuery{
		Username: username,
		Password: password,
	}
}

func (c AuthenticateUserQuery) Type() query.Type {
	return AuthenticateUserQueryType
}

type AuthenticateUserQueryHandler struct {
	service AuthenticateUserService
}

func NewAuthenticateUserQueryHandler(service AuthenticateUserService) AuthenticateUserQueryHandler {
	return AuthenticateUserQueryHandler{service: service}
}

func (h AuthenticateUserQueryHandler) Handle(ctx context.Context, query query.Query) response.Response {
	userQuery, ok := query.(AuthenticateUserQuery)
	if !ok {
		log.Printf("Invalid query type: %T\n", query)
		return response.NewResponse(http.StatusInternalServerError, "Internal Server Error")
	}
	return h.service.Authenticate(ctx, userQuery.Username, userQuery.Password)
}
