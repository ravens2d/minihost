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

	// h.templates["index.html"] = template.Must(template.ParseFiles("template/pages/index.tmpl", "template/base.tmpl"))
	// h.templates["register.html"] = template.Must(template.ParseFiles("template/pages/register.tmpl", "template/base.tmpl"))

	err := h.registerTemplates()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/login", h.Login)
	mux.HandleFunc("/logout", h.Logout)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", h.Home)

	h.muxHandler = h.session.LoadAndSave(mux)

	return h, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.muxHandler.ServeHTTP(w, r)
}
