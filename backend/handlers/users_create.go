package handlers

import (
	"net/http"
	"strings"
	"database/sql"
	"errors"
	"context"
	"log"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"

	"github.com/lib/pq"
)

type createUserParams struct {
	Email string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72,printascii"`
	Language database.Language `json:"language" validate:"omitempty,oneof=en es"`
	Role database.Role `json:"role" validate:"omitempty,oneof=visitor brand agency creator"`
	Name string `json:"name" validate:"omitempty,max=254"`
}

func (apiCfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := createUserParams{}
	err := decodeRequestBody(r, &params)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	createdUser, err := apiCfg.createUser(r.Context(), params)
	errCheckError := apiCfg.respondCreateUserError(w, err) 
	if errCheckError != nil {
		return
	}

	RespondWithJSON(w, http.StatusCreated, createUserResponse{
		User: createdUser,
	})
}

func (apiCfg *ApiConfig) respondCreateUserError(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrNotValidated):
		RespondWithError(w, http.StatusBadRequest, "Invalid email or password", err)
	case errors.Is(err, ErrHashingPassword):
		RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
	case errors.Is(err, ErrEmailTaken):
		RespondWithError(w, http.StatusConflict, "Email taken", err)
	case errors.Is(err, ErrCreatingUser):
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
	default:
	 	RespondWithError(w, http.StatusInternalServerError, "Unexpected error", err) 
	}
	return err
}

func (apiCfg *ApiConfig) createUser(ctx context.Context, params createUserParams) (User, error) {
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
		return User{}, ErrNotValidated
	}
 
	// hash password
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		return User{}, ErrHashingPassword
	}

	//create user
	createUserParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hash,
			Valid: true,
		},
		Language: params.Language,
		Role: params.Role,
		Name: sql.NullString{
			String: params.Name,
			Valid: params.Name != "",
		},
	}

	dbUser, err := apiCfg.DB.CreateUser(ctx, createUserParams)
	if err != nil {
		// check if the issue is because of a previous account with the same email (email is unique)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return User{}, ErrEmailTaken
		}
		log.Printf("Error creating user: %s", err)
		return User{}, ErrCreatingUser
	}

	// strip out the excess data
	handlersUser := dbUserToUser(dbUser)

	return handlersUser, nil
}
