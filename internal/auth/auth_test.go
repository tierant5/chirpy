package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/tierant5/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		want     string
		wantErr  bool
	}{
		{
			name:     "test1",
			password: "Pa$$word!",
			want:     "",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}
			if got == tt.password {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
			err := auth.CheckPasswordHash(tt.password, got)
			if err != nil {
				t.Errorf("CheckPasswordHash() = %v, failed unexpectedly", err)
			}
			err = auth.CheckPasswordHash("wrongPassword", got)
			if err == nil {
				t.Errorf("CheckPasswordHash() = %v, passed with incorrect password", err)
			}
			err = auth.CheckPasswordHash(tt.password, "")
			if err == nil {
				t.Errorf("CheckPasswordHash() = %v, passed with incorrect hash", err)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// Define test variables
	tokenSecret := "supersecretkey"
	userID := uuid.New()
	expiresIn := 2 * time.Hour

	// Generate a JWT
	token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("MakeJWT() failed to create token: %v", err)
	}

	// Validate the JWT
	parsedUserID, err := auth.ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("ValidateJWT() failed to parse token: %v", err)
	}
	if parsedUserID != userID {
		t.Errorf("ParsedUserID: %v != UserID: %v", parsedUserID, userID)
	}

	// Test with an invalid token
	_, err = auth.ValidateJWT("invalid.token.string", tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT should return an error for an invalid token")
	}

	// Test with an expired token
	expiredToken, _ := auth.MakeJWT(userID, tokenSecret, -1*time.Hour)
	_, err = auth.ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT should return an error for an expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			headers: http.Header{"Authorization": []string{"Bearer TOKEN_STRING"}},
			want:    "TOKEN_STRING",
			wantErr: false,
		},
		{
			name:    "test2",
			headers: http.Header{"NotAuthorization": []string{"Bearer TOKEN_STRING"}},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
