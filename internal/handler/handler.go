package handler

import (
	"minihost/internal/database"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// Handler ...
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	muxHandler http.Handler
	templates  *template.Template

	repo database.Repoistory
}

// New ...
func New(repo database.Repoistory) (Handler, error) {
	h := &handler{
		repo: repo,
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

	h.muxHandler = repo.SessionLoadAndSave(mux)

	return h, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.muxHandler.ServeHTTP(w, r)
}
