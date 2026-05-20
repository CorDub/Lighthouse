package handlers

import (
	"net/http"

	"Lighthouse/internal/database"

	"github.com/go-playground/validator/v10"
)

type ApiConfig struct {
	DB *database.Queries
	JWT string
	Env string
	Secure bool
	SameSite http.SameSite
	GoogleClientID string
	GoogleSecret string
	MailtrapUsername string
	MailtrapPassword string
	From string
	BaseURL string
}

var validate = validator.New(validator.WithRequiredStructEnabled())