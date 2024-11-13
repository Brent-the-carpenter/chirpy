package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}

func GetBearerToken(headers http.Header) (string, error) {
	if _, ok := headers["Authorization"]; !ok {
		return "", fmt.Errorf("Authorization header not present")
	}
	bearToken := headers.Get("Authorization")
	if bearToken == "" {
		return "", fmt.Errorf("Bearer token not set")
	}
	splitBearer := strings.Split(bearToken, " ")
	if len(splitBearer) != 2 {
		return "", fmt.Errorf("Improperly formatted bearer token")
	}
	return strings.TrimSpace(splitBearer[1]), nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("couldn't generate token: %v", err)
	}
	refreshToken := hex.EncodeToString(token)
	return refreshToken, nil
}
