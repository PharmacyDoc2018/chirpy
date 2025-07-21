package auth

import (
	"fmt"
	"testing"
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
