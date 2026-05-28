package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"database/sql"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email" validate:"required,email,max=254"`
		Password string `json:"password" validate:"required,min=8,max=72,printascii"`
	}

	//decode
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//normalize
	strings.TrimSpace(params.Email)
	strings.TrimSpace(params.Password)

	//validate
	errValidation := validate.Struct(params)
	if errValidation != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid email or password", errValidation)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	createUserParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hash,
			Valid: true,
		},
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, createUserResponse{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
}
