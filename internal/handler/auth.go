package handler

import (
	"minihost/internal/model"
	"minihost/internal/model/render"
	"minihost/internal/repository/database"
	"minihost/internal/util"
	"net/http"
	"net/mail"

	"go.uber.org/zap"
)

func (h *handler) RegisterGet(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := render.PopulateSessionInfo(r.Context(), h.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.RenderTemplate(w, "register.tmpl", render.PageState{SessionInfo: sessionInfo})
}

func (h *handler) LoginGet(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := render.PopulateSessionInfo(r.Context(), h.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.RenderTemplate(w, "login.tmpl", render.PageState{SessionInfo: sessionInfo})
}

func (h *handler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		h.RenderTemplate(w, "register.tmpl", render.PageState{
			ErrorMessage: "username, email, and password are required",
		})
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		h.RenderTemplate(w, "register.tmpl", render.PageState{
			ErrorMessage: "invalid email",
		})
		return
	}

	user, err := model.NewUser(username, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		util.Logger.Error(err)
		return
	}

	err = h.database.CreateUser(user)
	if err != nil {
		if err == database.ErrDuplicateUsername || err == database.ErrDuplicateEmail {
			h.RenderTemplate(w, "register.tmpl", render.PageState{
				ErrorMessage: err.Error(),
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			util.Logger.Error(err)
		}
		return
	}

	util.Logger.Infow("registered new user", zap.String("user_uuid", user.UUID.String()))
	h.session.SetAuthenticated(r.Context(), user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.database.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		util.Logger.Error(err)
		return
	}
	if user == nil {
		h.RenderTemplate(w, "login.tmpl", render.PageState{
			ErrorMessage: "account not found",
		})
		return
	}
	if !user.VerifyPassword(password) {
		h.RenderTemplate(w, "login.tmpl", render.PageState{
			ErrorMessage: "incorrect password",
		})
		return
	}

	h.session.SetAuthenticated(r.Context(), user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.session.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
