package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Function to generate a JWT token
func GenerateJWT(username string) (string, error) {

	var jwtSecret = []byte("my_secret_key")

	// Create a new token object with the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": "username",                           // Token subject (usually the user ID)
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Expiration time (1 hour)
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Token validation function
func validateBearerToken(r *http.Request) (User, error) {
	// Get the Authorization header from the HTTP request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header is missing")
		return User{}, fmt.Errorf("Authorization header is missing")
	}

	// Split the header to extract the token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Println("Invalid Authorization header format")
		return User{}, fmt.Errorf("Invalid Authorization header format")
	}

	// For simplicity, we're checking if the token is "valid-token" (replace with actual validation logic)
	token := parts[1]

	// ToDo: username error / unknown user not catched
	user, _ := getUserbyToken(token)
	if (user == User{}) {
		log.Println("Token for user not found - not authenticated")
		return User{}, fmt.Errorf("token not found")
	} else {
		//return "User: " + users.Username.(string), nil
		return user, nil
	}

}

func getcurrentTimestamp() string {
	// Define the layout
	layout := "02.01.2006 15:04"

	// Get the current date and time
	currentTime := time.Now()

	// Format the current time using the specified layout
	formattedTime := currentTime.Format(layout)
	return formattedTime
}
