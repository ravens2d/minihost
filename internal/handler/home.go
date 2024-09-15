package handler

import (
	"fmt"
	"io"
	"net/http"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	userUUID, err := h.repo.GetSessionAuthenticatedUserUUID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userUUID != nil {
		io.WriteString(w, fmt.Sprintf("logged in as %s", userUUID.String()))
		return
	}

	h.templates.ExecuteTemplate(w, "index.html", nil)
	return
}
