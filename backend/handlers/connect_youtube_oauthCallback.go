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

	//case where creator cancels or clicks no
	errParam := queryParams.Get("error")
	if errParam != "" {
		log.Printf("User denied OAuth consent: %s", errParam)

		// get oauth state
		state := queryParams.Get("state")
		fullOAuthState, lookupErr := apiCfg.DB.GetOAuthStateFromToken(r.Context(), state)
		if lookupErr != nil {
			apiCfg.sendOAuthPostMessage(w, "cancelled", "")
			return
		}

		// delete it
		_ = apiCfg.DB.DeleteOAuthStateFromToken(r.Context(), fullOAuthState.Token)
		apiCfg.sendOAuthPostMessage(w, "cancelled", fullOAuthState.ChannelHandle)
		return
	}

	// case where auth succeeds
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
		RespondWithError(w, 
			http.StatusBadGateway, 
			"Could not complete authorization", 
			errors.New("scope not present in token response."))
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

	// send a message to close the pop up and postMessage back to the original tab to update it
	apiCfg.sendOAuthPostMessage(w, "success", connection.ChannelHandle)
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
			log.Printf("OAuth state token not found: %v", err)
			return database.OauthState{}, ErrOAuthStateNotFound
		}

		log.Printf("Error getting the full oauth state from the db: %v", err)
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



func (apiCfg *ApiConfig) sendOAuthPostMessage(
w http.ResponseWriter, 
status string,
channelHandle string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<!DOCTYPE html>
		<html>
			<body>
				<script>
					if (window.opener) {
						window.opener.postMessage({
							source: "lighthouse-oauth",
							status: %q,
							channelHandle: %q
						}, %q);
					}
					window.close();
				</script>
			</body>
		</html>
	`, status, channelHandle, apiCfg.FrontEndOrigin)
}
