package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsconnections WSConnections

var dbEventBus *DBEventBus

// Serve the homepage (this serves the HTML page)
func serveHome(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "frontend/dist/index.html")

	// Serve the static folder at the root URL ("/")
	fs := http.FileServer(http.Dir("./frontend/dist/"))

	// Handle the root URL ("/") to serve files from the "./static" directory
	//http.Handle("/", fs)

	http.StripPrefix("/", fs).ServeHTTP(w, r)

	/*
		err := http.FileServer(http.Dir("frontend/dist/"))
		if err != nil {
			log.Fatalf("Error serving frontend: %s", err)
		} */
}

// Serve the homepage (this serves the HTML page)
func handleLogin(w http.ResponseWriter, r *http.Request) {

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the username and password
	var loginUser LoginUser
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate user credentials
	if !validateUser(loginUser.Username, loginUser.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token, err := GenerateJWT(loginUser.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// update token in DB

	errUpdToken := dbUpdateToken(loginUser.Username, token)
	if errUpdToken != nil {
		http.Error(w, "Failed to update token in db", http.StatusInternalServerError)
	}

	// Return the token in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// handleWebSocket function upgrades the HTTP connection to WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	// Validate the bearer token user
	user, err := validateBearerToken(r)
	if err != nil {
		log.Printf("user not found for token: %s", err)
		//http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		//return
	} else {
		log.Printf("connection from user: %s", user)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from WebSocket client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// Process the incoming message
		// extract only the token from header - full header is example [Bearer adminToken]
		response, err := processMessage(message, user)
		if err != nil {
			log.Println("Error processing message:", err)
			break
		}

		// Send the response back to the client
		err = conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}

	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", os.Getenv("DBFile"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dbEventBus = NewDBEventBus(db)

	errInit := initDB(db)
	if errInit != nil {
		log.Fatalf("Error init Database")
		os.Exit(1)
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not specified
	}

	http.HandleFunc("/", serveHome) // Serve the index page

	http.HandleFunc("/ws", handleWebSocket) // WebSocket endpoint

	http.HandleFunc("/login", handleLogin) // Login endpoint

	// Serve static files (e.g., JavaScript)
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//fmt.Printf("Server started on :%s\n", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))

	// Close the EventBus
	//eventBus.Close()

	// Start HTTP server
	server := &http.Server{Addr: ":" + port}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Start HTTP server
		log.Printf("Starting server on port :%s \n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("Shutting down server...")

	// Close the EventBus gracefully
	dbEventBus.Close()

	log.Println("Server stopped.")
}
