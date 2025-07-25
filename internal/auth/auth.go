package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject UUID: %w", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	const authPrefix = "Bearer "

	authValues := headers.Values("Authorization")
	if len(authValues) == 0 {
		return "", errors.New("error: authorization header not found")
	}

	authValue := authValues[0]
	if !strings.HasPrefix(authValue, authPrefix) {
		return "", errors.New("incorrect authorization header format")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authValue, authPrefix))
	if token == "" {
		return "", errors.New("authorization token is empty")
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	const tokenSize = 32

	token := make([]byte, tokenSize)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	hexToken := hex.EncodeToString(token)
	return hexToken, nil

}
