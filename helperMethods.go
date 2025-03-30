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

// Secret key used to sign and verify the JWT
var jwtSecret = []byte("your-secret-key")

// Custom claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Middleware to validate the Bearer token
func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Bearer token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Split the Bearer and token parts
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is what you expect (HMAC in this case)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		// Handle token parsing errors
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Ensure the token is valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check token expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		// Optionally, you can verify additional claims (e.g., check user roles, scopes)

		// Proceed with the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

// Token validation function
func validateBearerToken(token string) (User, error) {

	// ToDo: username error / unknown user not catched
	user, _ := getUserbyToken(token)
	if (user == User{}) {
		return User{}, fmt.Errorf("error: could not validate user - token not found: %s", token)
	} else {
		//return "User: " + users.Username.(string), nil
		return user, nil
	}

}

// Token validation function
func validateBearerToken_OLD(r *http.Request) (User, error) {
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
		return User{}, fmt.Errorf("ERROR: could not validate user - token not found: %s", token)
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

func checkDateTimeFormat(dateStr string) bool {

	layout := "02.01.2006 15:04"
	_, err := time.Parse(layout, dateStr)
	return err == nil // If parsing is successful, return true; otherwise, false
}

func calcDuration(fromStr string, toStr string) (string, error) {
	// calculate duration
	layout := "02.01.2006 15:04"
	startTime, err1 := time.Parse(layout, fromStr)
	endTime, err2 := time.Parse(layout, toStr)
	// Check for parsing errors
	if err1 != nil || err2 != nil {
		log.Printf("Error parsing dates: %s - %s", err1, err2)
		return "", fmt.Errorf("Error parsing dates for calcDuration")
	}
	// Calculate the duration between the two times
	timeDuration := endTime.Sub(startTime)
	// Convert the duration to hours and minutes
	totalHours := timeDuration.Hours()
	hours := int(totalHours)                           // Get the integer part (hours)
	minutes := int((totalHours - float64(hours)) * 60) // Get the fractional part as minutes

	// Format the duration as "hours.minutes"
	duration := fmt.Sprintf("%d.%02d", hours, minutes)

	return duration, nil

}
