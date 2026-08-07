package handlers

import (
	"net/http"
	"log"
	"context"
	"errors"
	"database/sql"
	"time"
	"fmt"

	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) YouTubeOAuthCallback(w http.ResponseWriter, r *http.Request) {
	// extract state and code query parameters
	queryParams := r.URL.Query()

	code := queryParams.Get("code")
	state := queryParams.Get("state")

	if code == "" || state == "" {
		log.Println("Missing code or state in OAuth callback")
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// check state is valid and not expired in oauth_states table in the DB
	fullOAuthState, err := apiCfg.checkOAuthStateValidity(r.Context(), state)

	// do the error handling
	errCheckError := apiCfg.respondOAuthStateError(w, err)
	if errCheckError != nil {
		return
	}

	// exchange code for tokens (refresh + access) with Google
	googleToken, err := apiCfg.YouTubeOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		RespondWithError(w, http.StatusBadGateway, "Could not complete authorization", err)
		return
	}

	// store the tokens in connections
	grantedScopes, ok := googleToken.Extra("scope").(string)
	if !ok {
		log.Println("Scopes granted unavailable through Extra()")
		RespondWithError(w, http.StatusBadGateway, "Could not complete authorization", err)
		return
	}
	connectionParams := database.CreateConnectionParams{
		UserID: fullOAuthState.UserID,
		Service: fullOAuthState.Service,
		ChannelID: fullOAuthState.ChannelID,
		ChannelHandle: fullOAuthState.ChannelHandle,
		AccessToken: googleToken.AccessToken,
		RefreshToken: sql.NullString{
			String: googleToken.RefreshToken,
			Valid: googleToken.RefreshToken != "",
		},
		TokenExpiresAt: sql.NullTime{
			Time: googleToken.Expiry,
			Valid: !googleToken.Expiry.IsZero(),
		},
		Scopes: grantedScopes,
	}

	connection, err := apiCfg.DB.CreateConnection(r.Context(), connectionParams)
	if err != nil {
		log.Printf("Could not create the connection: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Could not create the connection", err)
		return
	}

	// if valid proceed but first delete the oauth_state
	errDeletion := apiCfg.DB.DeleteOAuthStateFromToken(r.Context(), fullOAuthState.Token)
	if errDeletion != nil {
		log.Printf("Error deleting the oauth state: %s", errDeletion)
		RespondWithError(w, http.StatusInternalServerError, "Could not delete the OAuth state", err)
		return
	}

	// send a message to the waiting page to update it correctly
	

}


func (apiCfg *ApiConfig) respondOAuthStateError(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrOAuthStateNotFound):
		RespondWithError(w, http.StatusUnauthorized, "Invalid OAuth state", err)
	case errors.Is(err, ErrGettingOAuthState):
		RespondWithError(w, http.StatusInternalServerError, "Unexpected error", err)
	case errors.Is(err, ErrOAuthStateExpired):
		RespondWithError(w, http.StatusUnauthorized, "Invalid OAuth state", err)
	case errors.Is(err, ErrWrongOAuthService):
		RespondWithError(w, http.StatusUnauthorized, "Invalid OAuth state", err)
	default:
		RespondWithError(w, http.StatusInternalServerError, "Unexpected error", err)
	}
	return err
}



func (apiCfg *ApiConfig) checkOAuthStateValidity(ctx context.Context, token string) (database.OauthState, error) {
	//get full oauth state with token from db
	fullOAuthState, err := apiCfg.DB.GetOAuthStateFromToken(ctx, token)
	if err != nil {
		// no corresponding token / oauth state found
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("OAuth state token not found: %w", err)
			return database.OauthState{}, ErrOAuthStateNotFound
		}

		log.Printf("Error getting the full oauth state from the db: %w", err)
		return database.OauthState{}, ErrGettingOAuthState
	}

	//check expiry
	if fullOAuthState.ExpiresAt.Before(time.Now()) {
		log.Println("Expired oauth state token")
		return database.OauthState{}, ErrOAuthStateExpired
	}

	//check service
	if fullOAuthState.Service != "youtube" {
		log.Println("Wrong service")
		return database.OauthState{}, ErrWrongOAuthService
	}

	return fullOAuthState, nil
}