package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var token1, token2, token3 string

func TestJWTCreatingAndValidation(t *testing.T) {

	serverSecret := "animals"
	wrongSecret := "people"
	fiveSec := time.Second * 5
	oneSec := time.Second
	oneMinute := time.Minute
	test := []struct {
		userID     uuid.UUID
		JSONUserID string
		token      string
		deleteIn   time.Duration
		secret     string
	}{
		{
			JSONUserID: "chewy",
			token:      "",
			userID:     uuid.New(),
			deleteIn:   fiveSec,
			secret:     serverSecret,
		},
		{
			JSONUserID: "quartz",
			token:      "",
			userID:     uuid.New(),
			deleteIn:   oneMinute,
			secret:     wrongSecret,
		},
		{
			JSONUserID: "rocky",
			token:      "",
			userID:     uuid.New(),
			deleteIn:   oneSec,
			secret:     serverSecret,
		},
	}
	// Make tokens with different expiration times
	for index, tc := range test {
		token, err := MakeJWT(tc.userID, tc.secret, tc.deleteIn)
		if err != nil {
			t.Errorf("failed to make jwt, userID:%v , serverSecret: %s , duration: %v , error: %v",
				tc.userID,
				tc.secret,
				tc.deleteIn,
				err)
		}
		test[index].token = token
	}

	// Check that each token is able to be validated
	for _, tc := range test {
		_, err := ValidateJWT(tc.token, tc.secret)
		if err != nil {
			t.Errorf("failed to validate token: %s , secret: %s , error:%v",
				tc.token,
				tc.secret,
				err)
		}
	}

	// Make sure validation fails with wrong Secret key
	_, err := ValidateJWT(test[0].token, wrongSecret)
	if err == nil {
		t.Errorf("Validation should have failed. token secret: %s , provided secret: %s",
			test[0].secret,
			wrongSecret)
	}

	// Make sure tokens expire

	time.Sleep(time.Second * 3)
	_, err = ValidateJWT(test[2].token, serverSecret)
	if err == nil {
		t.Errorf("Validation should have failed, token expires:%v , slept for 3 seconds", test[2].deleteIn)
	}
}
