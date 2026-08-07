package middleware

import (
	"net/http"
	"context"
	
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
		userID, err := auth.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}) 
}