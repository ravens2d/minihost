package middleware

import (
	"minihost/internal/database"
	"net/http"
)

// RequireAuth ...
func RequireAuth(next http.HandlerFunc, repo database.Repoistory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID, err := repo.GetSessionAuthenticatedUserUUID(r.Context())
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
