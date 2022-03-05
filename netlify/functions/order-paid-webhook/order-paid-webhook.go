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

const TimestampFormat = "2006-01-02 15:04:05"

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if isNotAuthenticated(request) {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	_, err := createOrderInBioPlanet(request)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func isNotAuthenticated(request events.APIGatewayProxyRequest) bool {
	return checksum.CalculateWebhookChecksum(request, shoperApiKey) != request.Headers["x-webhook-sha1"]
}

func createOrderInBioPlanet(request events.APIGatewayProxyRequest) (*bioplanet.OrderConfirmation, error) {
	order, err := transformer.ToBioPlanetOrder(request)
	if err != nil {
		return nil, err
	}
	token, err := getBioPlanetApiToken()
	if err != nil {
		return nil, err
	}
	orderConfirmation, err := bioplanet.CreateOrder(*token, *order)
	if err != nil {
		fmt.Printf("Failed order creation: %+v\n", order)
		return nil, err
	}
	fmt.Printf("Order created no. %d\n", orderConfirmation.OrderId)
	return orderConfirmation, nil
}

func getBioPlanetApiToken() (*bioplanet.ApiToken, error) {
	utcTimeNow := time.Now().UTC().Format(TimestampFormat)
	apiTokenPost := bioplanet.ApiTokenPost{
		Hash:      checksum.CalculateTokenPostChecksum(bioPlanetApiKey, utcTimeNow, bioPlanetClientId),
		ClientId:  transformer.ToInt(bioPlanetClientId),
		Timestamp: utcTimeNow,
	}
	fmt.Println(apiTokenPost)
	return bioplanet.GetApiToken(apiTokenPost)
}

func main() {
	lambda.Start(handler)
}
