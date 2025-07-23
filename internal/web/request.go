package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayV2HTTPRequest

type Claims struct {
	Token string `json:"token"`
	Sub   string `json:"sub"`
}

func (r Request) GetClaims() (*Claims, error) {
	bearer, ok := r.Headers["Authorization"]
	if !ok {
		return nil, fmt.Errorf("missing Authorization header")
	}
	token := strings.Replace(bearer, "Bearer ", "", 1)
	segments := strings.Split(bearer, ".")
	if len(segments) != 3 {
		log.Printf("invalid token: %s\n", token)
		return nil, fmt.Errorf("invalid token segments: %d %+v", len(segments), segments)
	}

	// base64 decode the payload
	decodedBytes, err := base64.RawStdEncoding.DecodeString(strings.TrimSpace(segments[1]))
	if err != nil {
		return nil, fmt.Errorf("failed to decode the string: %v", err)
	}

	claims := &Claims{}

	err = json.Unmarshal(decodedBytes, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the string: %v", err)
	}

	claims.Token = token
	return claims, nil
}
