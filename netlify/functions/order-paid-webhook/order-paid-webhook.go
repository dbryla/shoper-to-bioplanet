package main

import (
	"crypto/sha1"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	hash := sha1.New()
	hash.Write([]byte(request.Headers["x-webhook-id"] + ":test-123:" + request.Body))
	checksum := string(hash.Sum(nil)[:])
	if checksum != request.Headers["x-webhook-sha1"] {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
