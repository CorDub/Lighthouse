package handlers

import (
	"net/http"
	"strings"
	"encoding/json"
	"database/sql"
	"errors"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"

	"github.com/lib/pq"
)

func (apiCfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email" validate:"required,email,max=254"`
		Password string `json:"password" validate:"required,min=8,max=72,printascii"`
		Language database.Language `json:"language" validate:"omitempty,oneof=en es"`
		Role database.Role `json:"role" validate:"omitempty,oneof=visitor brand agency creator"`
		Name string `json:"name" validate:"max=254"`
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
	params.Email = strings.TrimSpace(params.Email)
	params.Password = strings.TrimSpace(params.Password)
	params.Name = strings.TrimSpace(params.Name)

	//provide default values if nothing has been passed for optional fields
	if (params.Language == "") {
		params.Language = database.LanguageEN
	}
	if (params.Role == "") {
		params.Role = database.RoleVisitor
	}

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

	dbUser, err := apiCfg.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		// check if the issue is because of a previous account with the same email (email is unique)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			RespondWithError(w, http.StatusConflict, "Email taken", err)
			return
		}

		RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	handlersUser := dbUserToUser(dbUser)

	RespondWithJSON(w, http.StatusCreated, createUserResponse{
		User: handlersUser,
	})
}
