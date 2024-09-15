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

	repo, err := database.NewRepository()
	if err != nil {
		panic(err)
	}

	h, err := handler.New(repo)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8080", h)
	if err != nil {
		panic(err)
	}
}
