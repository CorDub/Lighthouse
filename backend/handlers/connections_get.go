package handlers

import (
	"net/http"
	"log"

	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) GetCreatorConnections(w http.ResponseWriter, r *http.Request) {
	//get user ID from middleware JWT validator
	userId, ok := GetUserIdFromJWT(r.Context())
	if !ok {
		log.Println("Unable to find the userID from the JWT")
		RespondWithError(w, http.StatusInternalServerError, "Unauthorized", nil)
		return
	}

	//get active connections
	activeConnections, err := apiCfg.DB.GetConnections(r.Context(), userId)
	if err != nil {
		log.Printf("Error getting connections: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Unexpected internal server error", err)
		return
	}

	//order connections by service with a map
	groupedByService := make(map[database.Service][]Connection)
	for _, connection := range activeConnections {
		//dbConnection to handlersConnection first, scour it for frontend
		cleanConnection := dbConnectionToConnection(connection)

		newServiceConnections := append(groupedByService[cleanConnection.Service], cleanConnection)
		groupedByService[cleanConnection.Service] = newServiceConnections
	}

	//send response
	RespondWithJSON(w, http.StatusOK, groupedByService)
}