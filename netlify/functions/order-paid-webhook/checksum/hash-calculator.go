package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/aws/aws-lambda-go/events"
	"hash"
	"strings"
	"time"
)

const TimestampFormat = "yyyy-MM-dd HH:mm:ss"

func CalculateWebhookChecksum(request events.APIGatewayProxyRequest, shoperApiKey string) string {
	sb := strings.Builder{}
	sb.WriteString(request.Headers["x-webhook-id"])
	sb.WriteString(":")
	sb.WriteString(shoperApiKey)
	sb.WriteString(":")
	sb.WriteString(request.Body)
	return calculateChecksum(sb.String(), sha1.New())
}

func CalculateTokenPostChecksum(bioPlanetApiKey string, now time.Time, bioPlanetClientId string) string {
	sb := strings.Builder{}
	sb.WriteString(bioPlanetApiKey)
	sb.WriteString(now.Format(TimestampFormat))
	sb.WriteString(bioPlanetClientId)
	return calculateChecksum(sb.String(), md5.New())
}

func calculateChecksum(data string, hash hash.Hash) string {
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)[:])
}
