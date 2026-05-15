package handlers

import (
	"net/http"
	"encoding/json"
	"time"
	"strings"
	"log"
	"database/sql"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email string `json:"email"`
}

type getUsersResponse struct {
	Users []User `json:"users"`
}

type createUserResponse struct {
		User
	}

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


func (apiCfg *ApiConfig) GetUsers(w http.ResponseWriter, r *http.Request) {
	//get users from DB
	dbUserList, err := apiCfg.DB.GetUsers(r.Context())
	if err != nil {
		log.Printf("Could not get the list of users: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "An unexpected server error occurred", err)
		return
	}

	// dbUser to handlersUser
	users := make([]User, len(dbUserList))
	for i, user := range dbUserList {
		users[i] = dbUserToUser(user)
	}

	//send answer
	RespondWithJSON(w, http.StatusOK, getUsersResponse{
		Users: users,
	})
}