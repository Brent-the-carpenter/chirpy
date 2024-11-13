package auth

import "testing"

func TestMakeRefreshToken(t *testing.T) {

	token, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("error occured while making token: %v", err)
	}
	if token == "" {
		t.Errorf("token is blank expected not blank string")
	}
}
