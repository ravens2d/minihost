package handler

import (
	"html/template"
	"net/http"

	"minihost/internal/repository/database"
	"minihost/internal/repository/session"
)

// Handler ...
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	muxHandler http.Handler
	templates  map[string]*template.Template

	database database.Database
	session  session.Session
}

// New ...
func New(db database.Database, s session.Session) (Handler, error) {
	h := &handler{
		templates: make(map[string]*template.Template),
		database:  db,
		session:   s,
	}

	err := h.registerTemplates()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /register", h.RegisterGet)
	mux.HandleFunc("GET /login", h.LoginGet)
	mux.HandleFunc("POST /register", h.RegisterPost)
	mux.HandleFunc("POST /login", h.LoginPost)
	mux.HandleFunc("GET /logout", h.Logout)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", h.Home)

	h.muxHandler = h.session.LoadAndSave(mux)

	return h, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.muxHandler.ServeHTTP(w, r)
}
