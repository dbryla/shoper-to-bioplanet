package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var apiKey = os.Getenv("API_KEY")

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if isNotAuthenticated(request) {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func isNotAuthenticated(request events.APIGatewayProxyRequest) bool {
	return request.QueryStringParameters["api-key"] != apiKey
}

func main() {
	lambda.Start(handler)
}
