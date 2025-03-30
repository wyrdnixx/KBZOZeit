package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"

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

// Server struct, you can add additional fields here as needed
type Server struct{}

// TemplateHandler to render HTML templates
func (s *Server) serveHTML(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Example data to pass to the HTML template
	data := struct {
		Title string
		Body  string
	}{
		Title: "Welcome to My Web Server",
		Body:  "This is a simple HTML page served by Go.",
	}

	// Render the template with data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}

var (
	// Global database connection
	DB *sql.DB

	// Mutex to ensure thread-safe initialization
	dbMutex sync.Mutex
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}

	// Connect to the SQLite database
	DB, err = sql.Open("sqlite3", os.Getenv("DBFile"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	dbEventBus = NewDBEventBus(DB)

	errInit := initDB(DB)
	if errInit != nil {
		log.Fatalf("Error init Database")
		os.Exit(1)
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if not specified
	}

	/*
		// Old test Handlers
		http.HandleFunc("/", indexHandler) // Serve the index page
		http.HandleFunc("/ws", handleWebSocket) // WebSocket endpoint
		http.HandleFunc("/login", loginHandler) // Login endpoint
		//http.HandleFunc("/home", homeHandler)   // Login endpoint
		http.Handle("/home", TokenValidationMiddleware(http.HandlerFunc(homeHandler)))
	*/

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/api/login", apiLoginHandler)
	http.HandleFunc("/app", appHandler)
	http.HandleFunc("/ws", wsHandler)

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
