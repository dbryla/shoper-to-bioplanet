package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/bioplanet"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/checksum"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/transformer"
	"os"
)

var shoperApiKey = os.Getenv("SHOPER_API_KEY")
var bioPlanetApiKey = os.Getenv("BIO_PLANET_API_KEY")
var bioPlanetClientId = os.Getenv("BIO_PLANET_CLIENT_ID")

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if isNotAuthenticated(request) {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	fmt.Println(request.Body)
	_, err := transformer.ToBioPlanetOrder(request)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	body, err := bioplanet.GetApiToken(bioPlanetApiKey, bioPlanetClientId)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	fmt.Println(string(body))

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func isNotAuthenticated(request events.APIGatewayProxyRequest) bool {
	return checksum.CalculateWebhookChecksum(request, shoperApiKey) != request.Headers["x-webhook-sha1"]
}

func main() {
	lambda.Start(handler)
}
