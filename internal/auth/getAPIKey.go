package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}

	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) != 2 || splitHeader[0] != "ApiKey" {
		return "", fmt.Errorf("Auth header is malformed")
	}
	return strings.TrimSpace(splitHeader[1]), nil

}
