package handlers

import (
	"net/http"
	"time"
	"log"
)

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