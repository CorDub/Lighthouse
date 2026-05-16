package handlers

import (
	"net/http"
	"net/url"

	"Lighthouse/internal/auth"
)

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
