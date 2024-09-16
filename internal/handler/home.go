package handler

import (
	"minihost/internal/model/render"
	"net/http"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // guard against catch all by default serve mux
		http.NotFound(w, r)
		return
	}

	sessionInfo, err := render.PopulateSessionInfo(r.Context(), h.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.RenderTemplate(w, "index.tmpl", sessionInfo)
	return
}
