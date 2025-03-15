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
func validateBearerToken(r *http.Request) (string, error) {
	// Get the Authorization header from the HTTP request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header is missing")
		return "", fmt.Errorf("Authorization header is missing")
	}

	// Split the header to extract the token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Println("Invalid Authorization header format")
		return "", fmt.Errorf("Invalid Authorization header format")
	}

	// For simplicity, we're checking if the token is "valid-token" (replace with actual validation logic)
	token := parts[1]

	// ToDo: username error / unknown user not catched
	users, _ := getUserbyToken(token)
	if (users == User{}) {
		log.Println("Token for user not found - not authenticated")
		return "", fmt.Errorf("token not found")
	} else {
		return "User: " + users.Username.(string), nil
	}
	//validToken := "asd3245345HDXKXKS3476hjdfksdfasdf"
	//validToken := "valid-token"
	//validToken := "adminToken" // ToDO - check against Database

	/* if token != validToken {
		log.Println("Invalid Bearer token")
		return "", fmt.Errorf("Invalid Bearer token")
	} */

	return "", fmt.Errorf("error fetching user from DB")

}
