package response

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

func HandleResponse(response Response) events.APIGatewayProxyResponse {
	bodyString, err := json.Marshal(response.GetValue())
	if err != nil {
		log.Println("Error marhsalling json:", err)
		return events.APIGatewayProxyResponse{
			Body:       `{"message" : "Internal Server Error"}"`,
			StatusCode: 500,
			Headers:    headers,
		}

	}

	return events.APIGatewayProxyResponse{
		Body:       string(bodyString),
		StatusCode: response.GetType(),
		Headers:    headers,
	}
}
