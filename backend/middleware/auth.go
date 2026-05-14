package middleware

import (
	"net/http"
	
	"Lighthouse/handlers"
	"Lighthouse/internal/auth"
)

func ValJWT(jwtSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get JWT token
		cookie, err := r.Cookie("token")
		if err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		tokenString := cookie.Value

		//validate JWT
		if _, err := auth.ValidateJWT(tokenString, jwtSecret); err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		next.ServeHTTP(w, r)
	}) 
}