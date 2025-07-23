package auth

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "123456"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error hashing password: %s", err)
	}
	fmt.Println(hashedPassword)

	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error password hash does not matched stored hash: %s", err)
	}

}

func TestMakeJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "123456"
	var expiresIn time.Duration = 5 * time.Second
	signedToken, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error creating signed token: %s", err)
	}
	fmt.Println(signedToken)

	userID, err := ValidateJWT(signedToken, tokenSecret)
	if err != nil {
		fmt.Println(err)
		t.Errorf("error validating token: %s", err)
	}

	fmt.Println(userID)
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		input         string
		expectedToken string
		expectedErr   error
	}{
		{
			input:         "Bearer 123456",
			expectedToken: "123456",
			expectedErr:   nil,
		},
		{
			input:         "123456",
			expectedToken: "",
			expectedErr:   errors.New("incorrect authorization header format"),
		},
		{
			input:         "",
			expectedToken: "",
			expectedErr:   errors.New("error: authorization headder not found"),
		},
	}

	for _, c := range cases {
		header := http.Header{}
		header.Set("Authorization", c.input)
		token, err := GetBearerToken(header)
		if token != c.expectedToken {
			t.Errorf("incorrect token. expected: %s. actual: %s", c.expectedToken, token)
		}
		if c.expectedErr == nil {
			if err != nil {
				t.Errorf("unexpected error. expected: nil. actual: %s", err)
			}
		} else if fmt.Sprint(err) != fmt.Sprint(c.expectedErr) {
			t.Errorf("unexpected error. expected: %s. actual: %s", c.expectedErr, err)
		}
	}
}
