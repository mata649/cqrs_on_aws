package adapter

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/mata649/cqrs_on_aws/user/internal/platform/server"
)

type LambdaAdapter struct {
	chiLambda *chiadapter.ChiLambda
}

func NewLambdaAdapter(server server.Server) *LambdaAdapter {
	return &LambdaAdapter{
		chiLambda: chiadapter.New(server.Engine()),
	}
}

func (l LambdaAdapter) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return l.chiLambda.ProxyWithContext(ctx, req)
}
