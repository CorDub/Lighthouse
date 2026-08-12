package handlers

import (
	"Lighthouse/internal/database"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) ToggleConnection(w http.ResponseWriter, r *http.Request) {
	// get userId from JWT
	userId, ok := GetUserIdFromJWT(r.Context())
	if !ok {
		log.Println("Unable to find the userID from the JWT")
		RespondWithError(w, http.StatusInternalServerError, "Unauthorized", nil)
		return
	}

	// get connection Id from params and convert it to uuid.UUID
	connectionIdStr := r.PathValue("id")
	connectionIdUUID, err := uuid.Parse(connectionIdStr)
	if err != nil {
		log.Printf("invalid string passed: %s", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// toggle active in DB
	updatedConnection, errDeactivating := apiCfg.toggleConnection(userId, connectionIdUUID, r.Context())
	if errDeactivating != nil {
		apiCfg.respondToggleConnectionError(w, errDeactivating)
		return
	}

	// respond
	RespondWithJSON(w, http.StatusOK, updatedConnection)
}


func (apiCfg *ApiConfig) toggleConnection(
	userId uuid.UUID, 
	connectionId uuid.UUID, 
	ctx context.Context) (Connection, error) {
	// get connection
	connection, err := apiCfg.DB.GetConnectionWithID(ctx, connectionId)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			log.Printf("No connection found: %s", err)
			return Connection{}, ErrConnectionNotFound
		}

		log.Printf("Error getting the connection: %s", err)
		return Connection{}, ErrGettingConnection
	}

	// check this is the correct user_id
	if connection.UserID != userId {
		log.Printf("Wrong userId asking for a modification of the connection: %s", err)
		return Connection{}, ErrUnauthorized
	}

	// flip to inactive
	params := database.ToggleConnectionParams{
		ID: connectionId,
		Active: !connection.Active,
	}
	updatedConnection, errDeactivating := apiCfg.DB.ToggleConnection(ctx, params)
	if err != nil {
		log.Printf("Error deactivating the connection: %s", errDeactivating)
		return Connection{}, ErrDeactivatingConnection
	}

	//scour connection
	cleanUpdatedConnection := dbConnectionToConnection(updatedConnection)

	return cleanUpdatedConnection, nil
}


func (apiCfg *ApiConfig) respondToggleConnectionError(w http.ResponseWriter, errorDeactivating error) {
	switch errorDeactivating {
	case ErrGettingConnection:
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", errorDeactivating)
	case ErrUnauthorized:
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", errorDeactivating)
	case ErrConnectionNotFound:
		RespondWithError(w, http.StatusBadRequest, "Connection not found", errorDeactivating)
	case ErrDeactivatingConnection:
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", errorDeactivating)
	}
} 