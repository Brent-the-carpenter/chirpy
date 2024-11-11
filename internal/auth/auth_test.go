package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "IamtheGreatestDev"
	pass2 := "Iamtheworstcatdad"

	hash, err := HashPassword(password)

	if password == hash {
		t.Errorf("Password and hash are the same. Expected to be different.")
	}

	err = CheckPasswordHash(password, hash)

	if err != nil {
		t.Errorf("Password did not produce hash. expected nil")
	}

	err = CheckPasswordHash(pass2, hash)
	if err == nil {
		t.Errorf("CheckPassword was supposed to return error, instead got nil")
	}

}
