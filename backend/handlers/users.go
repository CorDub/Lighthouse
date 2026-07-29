package handlers

import (
	"net/http"
	"time"
	"log"

	"Lighthouse/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email string `json:"email"`
	Language database.Language `json:"language"`
	Role database.Role `json:"role"`
	Name string `json:"name"`
}

type getUsersResponse struct {
	Users []User `json:"users"`
}

type createUserResponse struct {
		User
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