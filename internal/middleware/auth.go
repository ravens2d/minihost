package middleware

import (
	"minihost/internal/database"
	"net/http"
)

// RequireAuth ...
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID := database.SessionManager.GetString(r.Context(), database.UserUUIDSessionKey)
		if userUUID == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
