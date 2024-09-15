package handler

import (
	"fmt"
	"io"
	"minihost/internal/database"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	userUUID := database.SessionManager.GetString(r.Context(), database.UserUUIDSessionKey)

	if userUUID != "" {
		io.WriteString(w, fmt.Sprintf("logged in as %s", userUUID))
		return
	}

	// io.WriteString(w, "Welcome. Please log in")
	templates.ExecuteTemplate(w, "index.html", nil)
	return
}
