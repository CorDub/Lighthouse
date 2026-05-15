package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"net/url"
	"fmt"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"

	"google.golang.org/api/idtoken"
)

func (apiCfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	// decode params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// get user by email
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// return if no password (gogle login)
	if user.HashedPassword.Valid == false {
		RespondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	
	// check password is correct
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil || !match {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// get jwt token
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

	RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}



func (apiCfg *ApiConfig) Refresh(w http.ResponseWriter, r *http.Request) {
	//get refresh token
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	tokenString := cookie.Value

	// get token from db
	token, err := apiCfg.DB.GetRefreshTokenByToken(r.Context(), tokenString)

	// check if valid
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	if token.RevokedAt.Valid || token.ExpiresAt.Before(time.Now()) {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//if valid, revoke the old refresh token
	_, errRevoke := apiCfg.DB.RevokeRefreshToken(r.Context(), tokenString)
	if errRevoke != nil {
		log.Printf("Couldn't revoke the refresh token: %s", errRevoke)
		return
	}

	// create a new refresh token
	newRefreshToken := auth.MakeRefreshToken()
	newRefreshTokenParams := database.CreateRefreshTokenParams{
		Token: newRefreshToken,
		UserID: token.UserID,
	}
	
	if err := apiCfg.DB.CreateRefreshToken(r.Context(), newRefreshTokenParams); err != nil {
		log.Printf("Couldn't save the new refresh token in the DB : %s", err)
		return
	}

	// create also a new JWT cookie
	accessToken, err := auth.MakeJWT(token.UserID, apiCfg.JWT, time.Minute * 15)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// set the new refesh cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: newRefreshToken,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 30 * 24 * 60 * 60,
	})

	// set the new JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: accessToken,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 15 * 60,
	})

	RespondWithJSON(w, http.StatusOK, nil)
}



func (apiCfg *ApiConfig) Logout(w http.ResponseWriter, r *http.Request) {
	//get the refreshToken value
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	tokenString := cookie.Value

	//set refresh and JWT tokens to MaxAge -1 to clear them in the browser
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: "",
		HttpOnly: true,
		Expires: time.Unix(0, 0),
		Path: "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: "",
		HttpOnly: true,
		Expires: time.Unix(0, 0),
		Path: "/",
		MaxAge: -1,
	})

	// revoke the refresh token in the DB
	_, errRevoke := apiCfg.DB.RevokeRefreshToken(r.Context(), tokenString)
	if errRevoke != nil {
		log.Printf("Couldn't revoke the  refresh token: %s", errRevoke)
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", errRevoke)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}



func (apiCfg *ApiConfig) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// this is for the CSRF token, we can reuse the same function than for the refresh token
	token := auth.MakeRefreshToken()

	// set the CSRF token in http only cookie to keep it somewhere
	http.SetCookie(w, &http.Cookie{
		Name: "csrfToken",
		Value: token,
		Path: "/api/callback",
		HttpOnly: true,
		SameSite: apiCfg.SameSite,
		MaxAge: 300,
	})

	params := url.Values{
		"client_id": {apiCfg.GoogleClientID},
		"redirect_uri": {"http://localhost:8080/api/callback"},
		"response_type": {"code"},
		"scope": {"openid email"},
		"state": {token},
	}

	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()

	http.Redirect(w, r, authURL, http.StatusFound)
}



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
	if err != nil {
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


func (apiCfg *ApiConfig) CheckAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("check auth pinged")
	type respAuth struct {
		User User `json:"user"`
	}

	//get JWT token
	cookie, err := r.Cookie("token")

	// if you don't have the JWT check if the refresh token is present
	if err == http.ErrNoCookie {
		_, err := r.Cookie("refreshToken")
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		apiCfg.Refresh(w, r)
		return
	}
	
	// if it's there check the validity
	tokenString := cookie.Value
	userId, err := auth.ValidateJWT(tokenString, apiCfg.JWT)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//get the user
	dbUser, err := apiCfg.DB.GetUserByID(r.Context(), userId)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//prepare user to correct format for frontend
	handlersUser := dbUserToUser(dbUser)

	// prepare and send response
	userPayload := respAuth{
		User: handlersUser,
	}
	RespondWithJSON(w, http.StatusAccepted, userPayload)
}
