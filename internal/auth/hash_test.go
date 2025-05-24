package auth

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHashPassword(t *testing.T) {
	test_password := "hercules_12897"

	hashed_pwd, _ := HashPassword(test_password)

	if CheckPasswordHash(hashed_pwd, test_password) != nil {
		t.Errorf("Hashing Error")
	}

}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
