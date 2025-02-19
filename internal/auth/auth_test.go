package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.NewString()
	secret := "supersecret"
	tokenString, _ := MakeJWT(userID, secret, time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  string
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: tokenString,
			tokenSecret: secret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "wrong.token.string",
			tokenSecret: secret,
			wantUserID:  "",
			wantErr:     true,
		},
		{
			name:        "Invalid secret",
			tokenString: tokenString,
			tokenSecret: "wrongSecret",
			wantUserID:  "",
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(*testing.T) {
			uuid, err := ValidateJWT(test.tokenString, test.tokenSecret)
			if uuid != test.wantUserID {
				t.Errorf("ValidateJWT() uuid.UUID = %v, want %v", uuid, test.wantUserID)
			}
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateJWT() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tokenString, _ := MakeJWT(uuid.NewString(), "supersecret", time.Hour)
	bearer := fmt.Sprintf("Bearer %s", tokenString)
	headers := http.Header{
		"Authorization": []string{bearer},
	}

	tests := []struct {
		name            string
		headers         http.Header
		wantBearerToken string
		wantErr         bool
	}{
		{
			name:            "Bearer Token found",
			headers:         headers,
			wantBearerToken: tokenString,
			wantErr:         false,
		},
		{
			name:            "Authorization header not found",
			headers:         http.Header{},
			wantBearerToken: "",
			wantErr:         true,
		},
		{
			name:            "Bearer not found",
			headers:         http.Header{"Authorization": []string{}},
			wantBearerToken: "",
			wantErr:         true,
		},
		{
			name:            "Invalid Bearer",
			headers:         http.Header{"Authorization": []string{"TotallyWrongBearer token"}},
			wantBearerToken: "",
			wantErr:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenString, err := GetBearerToken(test.headers)
			if tokenString != test.wantBearerToken {
				t.Errorf("GetBearerToken() string = %v, want %v", tokenString, test.wantBearerToken)
			}
			if (err != nil) != test.wantErr {
				t.Errorf("GetBearerToken() err = %v, want %v", err, test.wantErr)
			}
		})
	}
}
