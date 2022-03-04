package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	hash := sha1.New()
	hash.Write([]byte(request.Headers["x-webhook-id"] + ":test-123:" + request.Body))
	fmt.Println("Expected " + request.Headers["x-webhook-sha1"])
	fmt.Printf("Actual: %x\n", hash.Sum(nil))
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
