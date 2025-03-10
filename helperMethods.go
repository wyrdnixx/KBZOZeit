package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Function to generate a JWT token
func GenerateJWT() (string, error) {

	var jwtSecret = []byte("my_secret_key")

	// Create a new token object with the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "1234567890",                         // Token subject (usually the user ID)
		"name": "John Doe",                           // Payload with any custom data
		"exp":  time.Now().Add(time.Hour * 1).Unix(), // Expiration time (1 hour)
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
