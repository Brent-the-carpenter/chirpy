package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	header := &http.Header{}
	header.Set("Authorization", "Bearer imLivingMyBestLife123231321")

	// Test getting bearer token
	token, err := GetBearerToken(*header)
	if err != nil {
		t.Errorf("error getting bearer token: %s, error:%v", header.Get("Authorization"), err)
	}

	if token != "imLivingMyBestLife123231321" {
		t.Errorf("Expected token: %s, got %s", header.Get("Authorization"), token)
	}

	//Test blank authorization
	header.Set("Authorization", "")
	if bearer, err := GetBearerToken(*header); err == nil {
		t.Errorf("expected to get bearer token not set error, got %s", bearer)
	}
	// Test improperly formatted bearer token
	header.Set("Authorization", "ILoveGolang")
	if bearer, err := GetBearerToken(*header); err == nil {
		t.Errorf("expected to get authorization token improperly formatted error, got %s", bearer)
	}
	// Test missing authorization header
	header.Del("Authorization")
	if bearer, err := GetBearerToken(*header); err == nil {
		t.Errorf("expected to get authorization header not set error, got %s", bearer)
	}

}
