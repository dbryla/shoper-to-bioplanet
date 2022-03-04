package bioplanet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/checksum"
	"github.com/dbryla/shoper-to-bioplanet/netlify/functions/order-paid-webhook/transformer"
	"io/ioutil"
	"net/http"
	"time"
)

func GetApiToken(bioPlanetApiKey string, bioPlanetClientId string) ([]byte, error) {
	utcTimeNow := time.Now().UTC()
	apiTokenPost, err := json.Marshal(ApiTokenPost{
		Hash:      checksum.CalculateTokenPostChecksum(bioPlanetApiKey, utcTimeNow, bioPlanetClientId),
		ClientId:  transformer.ToInt(bioPlanetClientId),
		Timestamp: utcTimeNow,
	})
	if err != nil {
		fmt.Println("Couldn't marshal bio planet api token.")
		return nil, err
	}
	response, err := http.Post("/api3/token", "application/json", bytes.NewBuffer(apiTokenPost))
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
