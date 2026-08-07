package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	//if channel exist, ask for OAuth
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



func (apiCfg *ApiConfig) prepareOAuthCall(ctx context.Context, handle string) (string, error) {
	//create state record in DB - userId, channelId, expiry

	//build consent URL

	//return it

	//frontend refirects browser with window.location.href = authURL

	// => creator approves

	// => returns from google at callback location

	// validate the state (record, expiry)

	// exchange code with google (POST to google's endpoint)

	// store the connection data in DB (connections table)

	// get back to app
}