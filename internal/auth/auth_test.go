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
	}
	fmt.Println(hashedPassword)

	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		fmt.Println(err)
	}

}

func TestMakeJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "123456"
	var expiresIn time.Duration = 5 * time.Second
	signedToken, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(signedToken)
}
