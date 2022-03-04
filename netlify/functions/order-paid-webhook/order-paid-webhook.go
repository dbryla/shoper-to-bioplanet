package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"strings"
)

var secret = os.Getenv("API_KEY")

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("Secret: %s\n", secret)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}
	checksum := calculateChecksum(request)
	fmt.Printf("Actual: %s\n", checksum)
	fmt.Printf("Expected: %s\n", request.Headers["x-webhook-sha1"])
	if checksum != request.Headers["x-webhook-sha1"] {
		fmt.Println("Rejected.")
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}
	fmt.Println("Authenticated.")
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func calculateChecksum(request events.APIGatewayProxyRequest) string {
	sb := strings.Builder{}
	sb.WriteString(request.Headers["x-webhook-id"])
	sb.WriteString(":")
	sb.WriteString(secret)
	sb.WriteString(":")
	sb.WriteString(request.Body)
	hash := sha1.New()
	hash.Write([]byte(sb.String()))
	checksum := string(hash.Sum(nil)[:])
	return checksum
}

func main() {
	lambda.Start(handler)
}
