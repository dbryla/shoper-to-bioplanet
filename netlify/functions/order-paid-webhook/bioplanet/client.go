package bioplanet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetApiToken(apiTokenPost ApiTokenPost) (*ApiToken, error) {
	request, err := json.Marshal(apiTokenPost)
	if err != nil {
		fmt.Println("Couldn't marshal bio planet api token request.")
		return nil, err
	}
	response, err := http.Post("https://drop.bioplanet.pl/api3/token", "application/json", bytes.NewBuffer(request))
	if err != nil {
		fmt.Println("Couldn't receive access token from bio planet.")
		return nil, err
	}
	defer response.Body.Close()
	return readTokenFromResponse(response)
}

func readTokenFromResponse(response *http.Response) (*ApiToken, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Couldn't read access token response from bio planet.")
		return nil, err
	}
	var apiToken ApiToken
	err = json.Unmarshal(body, &apiToken)
	if err != nil {
		fmt.Println("Couldn't unmarshal bio planet api token response.")
		return nil, err
	}
	return &apiToken, nil
}

func CreateOrder(token ApiToken, order Order) (*OrderConfirmation, error) {
	fmt.Printf("token %+v order %+v\n", token, order)
	requestBody, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Couldn't marshal bio planet order request.")
		return nil, err
	}
	response, err := sendPostToCreateOrder(token, requestBody)
	if err != nil {
		fmt.Println("Couldn't create order in bio planet.")
		return nil, err
	}
	defer response.Body.Close()
	return readOrderFromResponse(response)
}

func sendPostToCreateOrder(token ApiToken, requestBody []byte) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://drop.bioplanet.pl/api3/order", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Couldn't create request to create bio planet order.")
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	authorizationHeader := buildAuthorizationHeader(token)
	request.Header.Add("Authorization", authorizationHeader)
	return client.Do(request)
}

func readOrderFromResponse(response *http.Response) (*OrderConfirmation, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Couldn't read order response from bio planet.")
		return nil, err
	}
	var orderConfirmation OrderConfirmation
	err = json.Unmarshal(body, &orderConfirmation)
	if err != nil {
		fmt.Println("Couldn't unmarshal order confirmation response.")
		return nil, err
	}
	return &orderConfirmation, nil
}

func buildAuthorizationHeader(token ApiToken) string {
	sb := strings.Builder{}
	sb.WriteString(token.TokenType)
	sb.WriteString(" ")
	sb.WriteString(token.AccessToken)
	return sb.String()
}
