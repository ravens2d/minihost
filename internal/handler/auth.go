package handler

import (
	"minihost/internal/database"
	"minihost/internal/model"
	"minihost/internal/util"
	"net/http"
	"net/mail"

	"go.uber.org/zap"
)

// Register ...
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	user, err := model.NewUser(username, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		util.Logger.Error(err)
		return
	}

	// TODO: better error handling for duplcaite username or email
	err = database.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		util.Logger.Error(err)
		return
	}

	util.Logger.Infow("registered new user", zap.String("user_uuid", user.UUID.String()))
	database.SessionManager.Put(r.Context(), database.UserUUIDSessionKey, user.UUID.String())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := database.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		util.Logger.Error(err)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		util.Logger.Error(err)
		return
	}
	if !user.VerifyPassword(password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		util.Logger.Error(err)
		return
	}

	database.SessionManager.Put(r.Context(), database.UserUUIDSessionKey, user.UUID.String())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	database.SessionManager.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
