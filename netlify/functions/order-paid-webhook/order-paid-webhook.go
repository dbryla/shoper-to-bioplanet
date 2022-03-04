package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/bioplanet"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/checksum"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/transformer"
	"os"
	"time"
)

var shoperApiKey = os.Getenv("SHOPER_API_KEY")
var bioPlanetApiKey = os.Getenv("BIO_PLANET_API_KEY")
var bioPlanetClientId = os.Getenv("BIO_PLANET_CLIENT_ID")

const TimestampFormat = "yyyy-MM-dd HH:mm:ss"

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if isNotAuthenticated(request) {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	_, err := transformer.ToBioPlanetOrder(request)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	body, err := getBioPlanetApiToken()
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	fmt.Print(string(body))

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func getBioPlanetApiToken() ([]byte, error) {
	utcTimeNow := time.Now().UTC().Format(TimestampFormat)
	apiTokenPost := bioplanet.ApiTokenPost{
		Hash:      checksum.CalculateTokenPostChecksum(bioPlanetApiKey, utcTimeNow, bioPlanetClientId),
		ClientId:  transformer.ToInt(bioPlanetClientId),
		Timestamp: utcTimeNow,
	}
	fmt.Println(apiTokenPost)
	return bioplanet.GetApiToken(apiTokenPost)
}

func isNotAuthenticated(request events.APIGatewayProxyRequest) bool {
	return checksum.CalculateWebhookChecksum(request, shoperApiKey) != request.Headers["x-webhook-sha1"]
}

func main() {
	lambda.Start(handler)
}
