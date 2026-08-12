package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"context"

	"Lighthouse/internal/database"

	"github.com/google/uuid"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}


func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}


func dbUserToUser(user database.User) User {
	return User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Language: user.Language,
		Role: user.Role,
		Name: user.Name.String,
	}
}

func decodeRequestBody[T any](r *http.Request, payload *T) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(payload)
	if err != nil {
		return err
	}

	return nil
}

// this avoids clashes for third party libraries calling also a ctx.Value with string "userID"
// and doing so overwriting the result
// using a specific type ensures there's no collision
// function used in middleware - ValJWT but written here to avoid cyclical imports
type contextKey string
const UserIDKey contextKey = "userID"

func GetUserIdFromJWT(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func dbConnectionToConnection(connection database.Connection) Connection {
	return Connection{
		ID: connection.ID,
		Service: connection.Service,
		ChannelID: connection.ChannelID,
		ChannelHandle: connection.ChannelHandle,
		Scopes: connection.Scopes,
		Active: connection.Active,
	}
}