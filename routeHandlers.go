// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

var (
	templates *template.Template
)

func init() {
	/* 	// Create a test user with a hashed password
	   	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	   	users["admin"] = &User{
	   		Username: "admin",
	   		Password: string(hashedPassword),
	   	}
	*/
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
	var tokenString string

	// First try to get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// If not in header, check if this is a POST with form data (from login)
	if tokenString == "" && r.Method == "POST" {
		tokenString = r.FormValue("token")

		// If we got a token from POST, set it as a cookie for future requests
		if tokenString != "" {
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
			tokenString = cookie.Value
		}
	}

	// If no token found, redirect to login
	if tokenString == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate the token
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
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

	// For GET requests, serve the page
	if r.Method == "GET" {
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
func apiLoginHandler_OLD(w http.ResponseWriter, r *http.Request) {
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

	// Return token as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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
			conn.Close()
			return
		}

	}

}
