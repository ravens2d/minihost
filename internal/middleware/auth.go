package middleware

import (
	"minihost/internal/repository/session"
	"net/http"
)

// RequireAuth ...
func RequireAuth(next http.HandlerFunc, session session.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, err := session.GetAuthenticatedUserUUID(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if userUUID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
