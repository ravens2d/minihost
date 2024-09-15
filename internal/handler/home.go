package handler

import (
	"fmt"
	"io"
	"net/http"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	userUUID, err := h.session.GetAuthenticatedUserUUID(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userUUID != nil {
		io.WriteString(w, fmt.Sprintf("logged in as %s", userUUID.String()))
		return
	}

	h.RenderTemplate(w, "index.tmpl", nil)
	return
}
