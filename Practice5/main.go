package main

import (
	"log"
	"net/http"

	"practice5/db"
	"practice5/handler"
	"practice5/repository"
)

func main() {
	database := db.Connect()
	defer database.Close()

	repo := repository.New(database)
	h := handler.New(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", h.GetUsers)

	mux.HandleFunc("/users/common-friends", h.GetCommonFriends)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
