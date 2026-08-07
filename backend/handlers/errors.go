package handlers

import (
	"errors"
)

var (
	ErrNotValidated = errors.New("Invalid email or password")
	ErrHashingPassword = errors.New("Couldn't hash password")
	ErrEmailTaken = errors.New("Email taken")
	ErrCreatingUser = errors.New("Couldn't create user")
	ErrCreatingMagicLinkToken = errors.New("Couldn't create the magic link token")
	ErrOAuthStateNotFound = errors.New("Oauth state token not found")
	ErrGettingOAuthState = errors.New("Couldn't get OAuth state")
	ErrOAuthStateExpired = errors.New("OAuth state expired")
	ErrWrongOAuthService = errors.New("Wrong OAuth state service")
)