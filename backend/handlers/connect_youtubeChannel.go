package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"Lighthouse/internal/database"
	"Lighthouse/internal/auth"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (apiCfg *ApiConfig) ConnectYouTubeChannel (w http.ResponseWriter, r *http.Request) {
	// get the user ID from JWT
	userId, ok := GetUserIdFromJWT(r.Context())
	if !ok {
		log.Println("Unable to find the userID from the JWT")
		RespondWithError(w, http.StatusInternalServerError, "Unauthorized", nil)
		return
	}

	// decode body
	type connectYouTubeChannelBody struct {
		ChannelHandle string `json:"channelName" validate:"required,max=100,startswith=@"`
	}
	decodedBody := connectYouTubeChannelBody{}
	errDecoding := decodeRequestBody(r, &decodedBody) 
	if errDecoding != nil {
		log.Printf("Error decoding the body of the request: %s", errDecoding)
		RespondWithError(w, http.StatusInternalServerError, "Could not decode the JSON of the body", errDecoding)
		return
	}

	// normalize
	decodedBody.ChannelHandle = strings.TrimSpace(decodedBody.ChannelHandle)

	// validate
	errValidation := validate.Struct(decodedBody)
	if errValidation != nil {
		log.Printf("Error validating the body of the request: %s", errValidation)
		RespondWithError(w, http.StatusBadRequest, "The handle wasn't starting with @", errValidation)
		return
	}

	//verify channel exists
	channelID, err := apiCfg.verifyChannelExists(r.Context(), decodedBody.ChannelHandle)
	if err != nil {
		log.Printf("Error verifying channel exists: %s", err)
		RespondWithError(w, http.StatusBadGateway, "Could not verify channel", err)
		return
	}

	if channelID == "" {
		RespondWithError(w, http.StatusNotFound, "Channel not found", nil)
		return
	}

	//if channel exist, ask for Oauth
	// first create the ouath state token in db
	oauthStateToken, err := apiCfg.prepareOAuthState(r.Context(), userId, channelID, decodedBody.ChannelHandle)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not create the oauth state token", err)
		return
	}

	// build consent url
	oauthURL := apiCfg.buildConsentURL(oauthStateToken)

	// Respond with the url
	type ConnectYouTubeChannelResponseBody struct{
		AuthURL string `json:"authUrl"`
	}
	respBody := ConnectYouTubeChannelResponseBody{
		AuthURL: oauthURL,
	}
	RespondWithJSON(w, http.StatusOK, respBody)
}


func (apiCfg *ApiConfig) verifyChannelExists(ctx context.Context, handle string) (string, error) {
	// prepare call
	call := apiCfg.YouTubeService.Channels.List([]string{"id"}).ForHandle(handle)

	//call
	response, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("calling YouTube API: %w", err)
	}

	if len(response.Items) == 0 {
		return "", nil
	}

	return response.Items[0].Id, nil
}



func (apiCfg *ApiConfig) prepareOAuthState(
	ctx context.Context, 
	userId uuid.UUID,
	channelId string,
	channelHandle string,
	) (string, error) {
	//create state record in DB - userId, channelId, expiry
	oauthToken := auth.MakeRefreshToken()

	params := database.CreateOAuthStateParams{
		Token: oauthToken,
		UserID: userId,
		Service: database.ServiceYouTube,
		ChannelID: channelId,
		ChannelHandle: channelHandle,
	}

	oauthStateToken, err := apiCfg.DB.CreateOAuthState(ctx, params)
	if err != nil {
		log.Printf("Could not create the oauth state token: %s", err)
		return "", err
	}

	return oauthStateToken.Token, nil
}



func (apiCfg *ApiConfig) buildConsentURL(oauthStateToken string) string {
	authURL := apiCfg.YouTubeOAuthConfig.AuthCodeURL(
		oauthStateToken,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	return authURL
}

	// => returns from google at callback location

	// validate the state (record, expiry)

	// exchange code with google (POST to google's endpoint)

	// store the connection data in DB (connections table)

	// get back to app - in the same place with postMessage?