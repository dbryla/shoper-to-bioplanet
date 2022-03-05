package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/bioplanet"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/shoper"
	"strconv"
	"strings"
)

const ShoperCourierName = "Kurier"
const BioPlanetCourierName = "Kurier InPost"
const BioPlanetPaymentId = 85
const BioPlanetOrderComment = "Automatically created."

func ToBioPlanetOrder(request events.APIGatewayProxyRequest) (*bioplanet.Order, error) {
	var shoperOrder shoper.Order
	err := json.Unmarshal([]byte(request.Body), &shoperOrder)
	if err != nil {
		fmt.Println("Couldn't unmarshal shoper order.")
		return nil, err
	}

	var bioPlanetOrder = bioplanet.Order{
		Address: bioplanet.Address{
			Name:       buildName(shoperOrder),
			Street:     shoperOrder.BillingAddress.Street1,
			City:       shoperOrder.BillingAddress.City,
			PostalCode: shoperOrder.BillingAddress.Postcode,
			Phone:      shoperOrder.BillingAddress.Phone,
			Email:      shoperOrder.Email,
		},
		PaymentId:    BioPlanetPaymentId,
		DeliveryName: mapDeliveryName(shoperOrder),
		Comment:      BioPlanetOrderComment,
		OrderLines: bioplanet.OrderLines{
			KeyType: "Id",
		},
	}
	var lines []bioplanet.Line
	for _, product := range shoperOrder.Products {
		lines = append(lines, bioplanet.Line{Key: product.Code, Quantity: ToInt(product.Quantity)})
	}
	bioPlanetOrder.OrderLines.Lines = lines
	return &bioPlanetOrder, nil
}

func mapDeliveryName(shoperOrder shoper.Order) string {
	if ShoperCourierName == shoperOrder.Shipping.Name {
		return BioPlanetCourierName
	}
	return BioPlanetCourierName
}

func buildName(shoperOrder shoper.Order) string {
	sb := strings.Builder{}
	sb.WriteString(shoperOrder.BillingAddress.Firstname)
	sb.WriteString(" ")
	sb.WriteString(shoperOrder.BillingAddress.Lastname)
	return sb.String()
}

func ToInt(stringInt string) int {
	result, _ := strconv.Atoi(stringInt)
	return result
}
