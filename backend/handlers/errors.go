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
)