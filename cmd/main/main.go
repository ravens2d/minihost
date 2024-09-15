package main

import (
	"fmt"
	"net/http"

	"minihost/internal/handler"
	"minihost/internal/repository/database"
	"minihost/internal/repository/session"
)

func main() {
	fmt.Println("running...")

	db, err := database.New()
	if err != nil {
		panic(err)
	}
	session, err := session.New()
	if err != nil {
		panic(err)
	}

	h, err := handler.New(db, session)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8080", h)
	if err != nil {
		panic(err)
	}
}
