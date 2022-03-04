package main

import (
	"encoding/json"
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

	body, err := getBioPlanetApiToken(err)
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

func getBioPlanetApiToken(err error) ([]byte, error) {
	utcTimeNow := time.Now().UTC()
	apiTokenPost, err := json.Marshal(bioplanet.ApiTokenPost{
		Hash:      checksum.CalculateTokenPostChecksum(bioPlanetApiKey, utcTimeNow, bioPlanetClientId),
		ClientId:  transformer.ToInt(bioPlanetClientId),
		Timestamp: utcTimeNow,
	})
	if err != nil {
		fmt.Println("Couldn't marshal bio planet api token.")
		return nil, err
	}
	body, err := bioplanet.GetApiToken(apiTokenPost)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func isNotAuthenticated(request events.APIGatewayProxyRequest) bool {
	return checksum.CalculateWebhookChecksum(request, shoperApiKey) != request.Headers["x-webhook-sha1"]
}

func main() {
	lambda.Start(handler)
}
