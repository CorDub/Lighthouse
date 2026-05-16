package handlers

import (
	"net/http"
)

func (apiCfg *ApiConfig) checkPassword(w http.ResponseWriter, r *http.Request) {
	//extract email sent
	emailSent := r.Body.Get("email")

	//validate email sent
}