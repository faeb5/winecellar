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
)

const tokenIssuer = "secretsanta"

func MakeJWT(userID string, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID,
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)

	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", fmt.Errorf("Unable to get issuer from claims: %s", err)
	}
	if issuer != tokenIssuer {
		return "", errors.New("Invalid issuer")
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", fmt.Errorf("Unable to get expiration time from claims: %s", err)
	}
	if expirationTime.UTC().Before(time.Now().UTC()) {
		return "", errors.New("Token is expired")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve subject from claims: %s", err)
	}

	return subject, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	return getValueFromAuthHeader(headers, "Bearer")
}

func getValueFromAuthHeader(headers http.Header, key string) (string, error) {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("Authorization header not found")
	}

	var keyValue string
	for _, val := range authHeader {
		fields := strings.Fields(val)
		if len(fields) != 2 {
			continue
		}
		if strings.EqualFold(fields[0], key) {
			keyValue = fields[1]
		}
	}
	if keyValue == "" {
		return "", fmt.Errorf("No value for key %s found", key)
	}

	return keyValue, nil
}

func MakeRefreshToken() (string, error) {
	rawData := make([]byte, 32)
	if _, err := rand.Read(rawData); err != nil {
		return "", errors.New("Unable to create refresh token")
	}
	return hex.EncodeToString(rawData), nil
}
