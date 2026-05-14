package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}


func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "lighthouse",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey := []byte(tokenSecret)

	signedString, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("Could not sign the jwt token: %s", err)
		return "", err
	}

	return signedString, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString, 
		&claimsStruct, 
		func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil })
	if err != nil {
		log.Printf("Could not parse the token provided: %s", err)
		return uuid.Nil, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claimsStruct.Subject)
	if err != nil {
		log.Printf("Couldn't parse the user UUID from the claims: %s", err)
		return uuid.Nil, err
	}

	return userID, nil
}


func GetBearerToken(headers http.Header) (string, error) {
	fullBearerToken := headers.Get("Authorization")
	if fullBearerToken == "" {
		log.Println("No bearer token found")
		return "", errors.New("No bearer token found")
	}

	if !strings.HasPrefix(fullBearerToken, "Bearer ") {
		log.Println("Malformed token")
		return "", errors.New("Malformed token")
	}

	strippedBearerToken := strings.TrimPrefix(fullBearerToken, "Bearer ")
	return strippedBearerToken, nil
}


func MakeRefreshToken() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	
	token := hex.EncodeToString(buffer)
	return token
}