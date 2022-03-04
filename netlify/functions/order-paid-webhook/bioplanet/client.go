package bioplanet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetApiToken(apiTokenPost ApiTokenPost) ([]byte, error) {
	request, err := json.Marshal(apiTokenPost)
	if err != nil {
		fmt.Println("Couldn't marshal bio planet api token.")
		return nil, err
	}
	response, err := http.Post("https://drop.bioplanet.pl/api3/token", "application/json", bytes.NewBuffer(request))
	if err != nil {
		fmt.Println("Couldn't receive access token from bio planet.")
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Couldn't read access token from bio planet.")
		return nil, err
	}
	return body, nil
}
