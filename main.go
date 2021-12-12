package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/wyrdnixx/KBZOZeit/frontend/api"
	"github.com/wyrdnixx/KBZOZeit/frontend/database"
)

func main() {
	database.Initdb()
	startWebServer()
}

func startWebServer() {
	// Choose the folder to serve
	staticDir := "/frontend/dist"

	router := mux.NewRouter()

	router.HandleFunc("/api/TestApi", api.TestApi)
	// Create the route
	router.
		PathPrefix("/").
		Handler(http.StripPrefix("/", http.FileServer(http.Dir("."+staticDir))))

	srv := &http.Server{
		Handler: router,
		Addr:    ":8081",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
