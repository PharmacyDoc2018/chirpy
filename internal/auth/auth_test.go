package auth

import (
	"fmt"
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
