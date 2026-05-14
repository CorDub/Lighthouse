package handlers

import (
	"Lighthouse/internal/database"

	"github.com/go-playground/validator/v10"
)

type ApiConfig struct {
	DB *database.Queries
	JWT string
	Env string
}

var validate = validator.New(validator.WithRequiredStructEnabled())