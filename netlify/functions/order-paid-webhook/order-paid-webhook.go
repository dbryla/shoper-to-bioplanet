package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/model"
	"os"
	"strconv"
	"strings"
)

var secret = os.Getenv("API_KEY")

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if isNotAuthenticated(request) {
		return &events.APIGatewayProxyResponse{
			StatusCode: 401,
		}, nil
	}

	var shoperOrder model.ShoperOrder
	err := json.Unmarshal([]byte(request.Body), &shoperOrder)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}
	fmt.Println(shoperOrder)

	//var bioPlanetOrder = model.BioPlanetOrder{
	//	Address: model.Address{
	//		Name:       buildName(shoperOrder),
	//		Street:     shoperOrder.BillingAddress.Street1,
	//		City:       shoperOrder.BillingAddress.City,
	//		PostalCode: shoperOrder.BillingAddress.Postcode,
	//		Phone:      shoperOrder.BillingAddress.Phone,
	//		Email:      shoperOrder.Email,
	//	},
	//	PaymentId:    toInt(shoperOrder.PaymentId),
	//	DeliveryName: shoperOrder.Shipping.Name,
	//	Comment:      "Automatically created.",
	//	OrderLines: struct {
	//		KeyType string `json:"KeyType"`
	//		Lines   []struct {
	//			Key      string `json:"Key"`
	//			Quantity int    `json:"Quantity"`
	//		} `json:"Lines"`
	//	}{},
	//	InpostPaczkomatCode: "",
	//}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func toInt(stringInt string) int {
	result, _ := strconv.Atoi(stringInt)
	return result
}

func buildName(shoperOrder model.ShoperOrder) string {
	sb := strings.Builder{}
	sb.WriteString(shoperOrder.BillingAddress.Firstname)
	sb.WriteString("")
	sb.WriteString(shoperOrder.BillingAddress.Lastname)
	return sb.String()
}

func isNotAuthenticated(request events.APIGatewayProxyRequest) bool {
	return calculateChecksum(request) != request.Headers["x-webhook-sha1"]
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
	checksum := hex.EncodeToString(hash.Sum(nil)[:])
	return checksum
}

func main() {
	lambda.Start(handler)
}
