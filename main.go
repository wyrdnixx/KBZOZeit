package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/wyrdnixx/KBZOZeit/api"
	"github.com/wyrdnixx/KBZOZeit/database"
	"github.com/wyrdnixx/KBZOZeit/models"
)

var Cfg models.Configuration

func main() {

	database.Initdb()
	startWebServer()
}

func startWebServer() {
	// Choose the folder to serve
	staticDir := "/frontend/dist"

	router := mux.NewRouter()

	router.HandleFunc("/api/TestApi", api.TestApi)
	router.HandleFunc("/api/AdminGetUsers", api.AdminGetUsers)
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
