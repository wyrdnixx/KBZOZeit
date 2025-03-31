// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"

	"golang.org/x/crypto/bcrypt"
)

var (
	templates *template.Template
)

func init() {

	// Parse templates
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func appHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s : /app requested", r.RemoteAddr)

	var tokenString string

	// First try to get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
		log.Printf("%s : Got Baerer token ", r.RemoteAddr)
	}

	// If not in header, check if this is a POST with form data (from login)
	if tokenString == "" && r.Method == "POST" {
		tokenString = r.FormValue("auth_token")

		// If we got a token from POST, set it as a cookie for future requests
		if tokenString != "" {
			log.Printf("%s : Got login token - setting cookie", r.RemoteAddr)
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    tokenString,
				Path:     "/",
				HttpOnly: true,
				Secure:   r.TLS != nil,
				MaxAge:   86400, // 24 hours
			})
		}
	}

	// If still no token, check for cookie
	if tokenString == "" {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			log.Printf("%s : got auth_token", r.RemoteAddr)
			tokenString = cookie.Value
		}
	}

	// If no token found, redirect to login
	if tokenString == "" {
		log.Printf("%s : no token found - redirect ro login ", r.RemoteAddr)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate the token
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {

		log.Printf("%s : Clear invalid token", r.RemoteAddr)
		// Clear invalid cookie if it exists
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1, // Delete the cookie
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// First try to authenticate with cookie
	var username string

	if err == nil && token.Valid {
		username = claims.Subject
		log.Printf("%s : cookie token authenticated successfully: %s", r.RemoteAddr, username)
	}

	// For GET requests, serve the page
	if r.Method == "GET" {
		log.Printf("%s : %s : User authenticated - display /app", r.RemoteAddr, username)
		templates.ExecuteTemplate(w, "app.html", map[string]string{
			"Username": claims.Subject,
		})
		return
	}

	// For POST requests (from login), redirect to GET to prevent form resubmission
	http.Redirect(w, r, "/app", http.StatusSeeOther)
}
func apiLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	userStoredPasswd, err := getUserPasswordHash(username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)

	}

	if err := bcrypt.CompareHashAndPassword([]byte(userStoredPasswd), []byte(password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	claims := jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Set the token as a cookie (will be used as a fallback)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   86400, // 24 hours
	})

	// Return token as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s : User logout requested. deleting cookie", r.RemoteAddr)
	// Clear session cookie
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Expire it immediately
		HttpOnly: true,
		Secure:   true, // Set to true if you're using HTTPS
	}
	http.SetCookie(w, cookie)

	// Additional logic to clear session or tokens if needed

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s : new websocket connection request.", r.RemoteAddr)
	// Upgrade connection first, then handle authentication with the first message
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// First try to authenticate with cookie
	var username string
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		// Validate token from cookie
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err == nil && token.Valid {
			username = claims.Subject
			log.Printf("%s : cookie token authenticated successfully: %s", r.RemoteAddr, username)
		}
	}

	// If cookie authentication failed, wait for auth message
	if username == "" {
		// Set a timeout for authentication
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// Wait for authentication message
		var authMsg struct {
			Type  string `json:"type"`
			Token string `json:"token"`
		}

		if err := conn.ReadJSON(&authMsg); err != nil || authMsg.Type != "auth" {
			log.Printf("%s : no or wrong authentication message: %s", r.RemoteAddr, authMsg.Type)
			conn.Close()
			return
		}

		// Reset the deadline
		conn.SetReadDeadline(time.Time{})

		// Validate token
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(authMsg.Token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Printf("%s : token via auth_message authentication error: %s", r.RemoteAddr, err)
			conn.Close()
			return
		}

	}

	for {
		// Read message from WebSocket client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// Process the incoming message

		user, err := getUserbyName((username))
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("error - user not found."+err.Error()))
		} else {
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

}
