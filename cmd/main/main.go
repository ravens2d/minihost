package main

import (
	"fmt"
	"net/http"

	"minihost/internal/database"
	_ "minihost/internal/database"
	"minihost/internal/handler"
)

func main() {
	fmt.Println("running...")

	mux := http.NewServeMux()

	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/logout", handler.Logout)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handler.Home)

	err := http.ListenAndServe(":8080", database.SessionManager.LoadAndSave(mux))
	if err != nil {
		panic(err)
	}
}
