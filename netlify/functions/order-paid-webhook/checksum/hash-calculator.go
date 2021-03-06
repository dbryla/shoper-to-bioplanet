package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/aws/aws-lambda-go/events"
	"hash"
	"strings"
)

func CalculateWebhookChecksum(request events.APIGatewayProxyRequest, shoperApiKey string) string {
	sb := strings.Builder{}
	sb.WriteString(request.Headers["x-webhook-id"])
	sb.WriteString(":")
	sb.WriteString(shoperApiKey)
	sb.WriteString(":")
	sb.WriteString(request.Body)
	return calculateChecksum(sb.String(), sha1.New())
}

func CalculateTokenPostChecksum(bioPlanetApiKey string, timestamp string, bioPlanetClientId string) string {
	sb := strings.Builder{}
	sb.WriteString(bioPlanetApiKey)
	sb.WriteString(timestamp)
	sb.WriteString(bioPlanetClientId)
	return calculateChecksum(sb.String(), md5.New())
}

func calculateChecksum(data string, hash hash.Hash) string {
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)[:])
}
