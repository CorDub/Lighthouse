package auth_test

import (
	"testing"
	"time"
	"net/http"

	"Lighthouse/internal/auth"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := auth.HashPassword(password1)
	hash2, _ := auth.HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := auth.CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}


func TestMakeJWT(t *testing.T) {

	tests := []struct {
		name string
		userId uuid.UUID
		tokenSecret string
		expiresIn time.Duration
		wantErr bool
	}{
		{
			name: "correct",
			userId: uuid.New(),
			tokenSecret: "r1JBHToaXvmgYGobtygRFNhlssNBv7VFwTqNeh5Tq48=",
			expiresIn: 2 * time.Hour,
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// check error presence
			_, err := auth.MakeJWT(test.userId, test.tokenSecret, test.expiresIn)
			if err != nil {
				if !test.wantErr {
					t.Errorf("Error: %s", err)
					return
				}
			}
		})
	}
}


func TestValidateJWT(t *testing.T) {
	 tests := []struct {
		name string
		issuer string
		userID uuid.UUID
		expiresIn time.Duration
		tokenSecret string
		wantErr bool
	 }{
		{
			name: "correct validation",
			userID: uuid.New(),
			expiresIn: time.Hour * 2,
			tokenSecret: "r1JBHToaXvmgYGobtygRFNhlssNBv7VFwTqNeh5Tq48=",
			wantErr: false,
		},
		{
			name: "expired token",
			issuer: "lighthouse",
			userID: uuid.New(),
			expiresIn: time.Nanosecond,
			tokenSecret: "r1JBHToaXvmgYGobtygRFNhlssNBv7VFwTqNeh5Tq48=",
			wantErr: true,
		},
		{
			name: "wrong token secret",
			issuer: "lighthouse",
			userID: uuid.New(),
			expiresIn: time.Hour * 2,
			tokenSecret: "EcKBwFKyCSN2m0jJYb+mmnkkuItONA5gQoBTAIp+jzA=",
			wantErr: true,
		},
	 }

	 for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create token based on what is passed
			tokenString, err := auth.MakeJWT(test.userID, test.tokenSecret, test.expiresIn)
			if err != nil {
				t.Fatalf("Could not create the JWT: %s", err)
			}
			
			// check validity
			_, err2 := auth.ValidateJWT(tokenString, "r1JBHToaXvmgYGobtygRFNhlssNBv7VFwTqNeh5Tq48=")
			if err2 != nil {
				if !test.wantErr {
					t.Errorf("Unexpected error: %s", err2)
					return
				}
			}
		})
	 }
}


func TestGetBearerToken(t *testing.T) {
	token1 := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	strippedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	malformedToken := "InR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	tests :=  []struct {
		name string
		token string
		wantErr bool
	}{
		{
			name: "correct",
			token: token1,
			wantErr: false,
		},
		{
			name: "no authorization header",
			token: token1,
			wantErr: true,
		},
		{
			name: "malformed token",
			token: malformedToken,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create header based on what is passed
			headers := http.Header{}
			if test.name != "no authorization header" {
				headers.Add("Authorization", test.token)
			}

			// check error presence
			bearerToken, err := auth.GetBearerToken(headers)

			if !test.wantErr {
				if err != nil {
					t.Errorf("Unexpected error: %s", err)
					return
				} else {
					
					// check correct return value
					if bearerToken != strippedToken {
						t.Errorf("Got: %s, Expected: %s", bearerToken, strippedToken)
						return
					}
				}
			} else {
				if err == nil {
					t.Errorf("Expected error but got nothing")
					return
				}
			}
		})
	}
}