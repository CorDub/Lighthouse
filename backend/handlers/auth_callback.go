package handlers

import (
	"net/http"
	"net/url"
	"encoding/json"
	"database/sql"
	"time"
	"log"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"

	"google.golang.org/api/idtoken"
)

func (apiCfg *ApiConfig) Callback(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	// get CSRF token from cookie
	cookie, err := r.Cookie("csrfToken")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	tokenString := cookie.Value

	// compare it to the state sent back from Google
	stateFromGoogle := r.URL.Query().Get("state")
	match := tokenString == stateFromGoogle

	// delete CSRF cookie
	http.SetCookie(w, &http.Cookie{
		Name: "csrfToken",
		Value: "",
		Path: "/api/callback",
		HttpOnly: true,
		SameSite: apiCfg.SameSite,
		MaxAge: -1,
	})

	if !match {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// get the code google sends from the params
	code := r.URL.Query().Get("code")

	// now send it back to get confirmation on their end and all the data needed
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
	"code": {code},
	"client_id": {apiCfg.GoogleClientID},
	"client_secret": {apiCfg.GoogleSecret},
	"redirect_uri": {"http://localhost:8080/api/callback"},
	"grant_type": {"authorization_code"},
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to exchange code", err)
		return
	}
	defer resp.Body.Close()

	//decode the response
	var tokenResponse struct{
		IDToken string `json:"id_token"`
	}
	errDecoding := json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if errDecoding != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode Google's answer", errDecoding)
		return
	}

	//verify Google signature
	payload, err := idtoken.Validate(r.Context(), tokenResponse.IDToken, apiCfg.GoogleClientID)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//get the email from payload
	email := payload.Claims["email"].(string)

	//check if user exist in the database
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), email)

	//create if doesn't exist
	if err == sql.ErrNoRows {
		createdUser, err := apiCfg.DB.CreateUserViaGoogle(r.Context(), email)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Couldn't create a new user via Google", err)
			return
		}
		user = createdUser
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	//log in the user - make the jwt
	token, err := auth.MakeJWT(user.ID, apiCfg.JWT, time.Minute * 15)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// set the JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: token,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 15 * 60,
	})

	// create refresh token
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	}

	// save it in DB
	if err := apiCfg.DB.CreateRefreshToken(r.Context(), refreshTokenParams); err != nil {
		log.Printf("Couldn't save the refresh token in db: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// set the refresh cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 30 * 24 * 60 * 60,
	})

	http.Redirect(w, r, "http://localhost:5173/home", http.StatusFound)
}