package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"minihost/internal/repository/database"
	"minihost/internal/repository/session"
)

// Handler ...
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	muxHandler http.Handler
	templates  *template.Template

	database database.Database
	session  session.Session
}

// New ...
func New(db database.Database, s session.Session) (Handler, error) {
	h := &handler{
		database: db,
		session:  s,
	}

	h.templates = template.New("")
	err := filepath.Walk("template", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = h.templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
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
