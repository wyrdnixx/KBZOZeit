package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// handleEchoTest processes echoTest type messages
func handleEchoTest(content interface{}) ([]byte, error) {
	// Convert content to string (it's expected to be a string)
	contentStr, ok := content.(string)
	if !ok {
		return generateResponse("handleEchoTestResponse", true, "Invalid content for echoTest")
	}

	// Respond with the same content
	response := map[string]string{
		"type":    "echoResponse",
		"content": contentStr,
	}
	return json.Marshal(response)
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
	if !validateUser(loginUser.Username, loginUser.PwdHash) {
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
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
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

// handleTimeBooking processes timebooking type messages
func handleTimeBooking(content interface{}, user User) ([]byte, error) {
	// First, convert content to a map
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return generateResponse("handleTimeBookingResponse", true, "Invalid content format for timebooking")
	}

	// Extract the "from" and "to" fields
	fromStr, okFrom := contentMap["from"].(string)
	toStr, okTo := contentMap["to"].(string)
	/* 	if !okFrom || !okTo {
		return generateResponse("handleTimeBookingResponse", true, "Missing 'from' or 'to' in timebooking content")
	} */

	if okFrom && !okTo { // only clock in

		// check for valid dateTime Format
		if !checkDateTimeFormat(fromStr) {
			log.Printf("ERROR: From string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: From string not in valid DateTime format.")
		}
		// insert from to new booking
		log.Printf(`got "from" booking: %s`, fromStr)

		// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if res { // Error - user has open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has already open booking")
		}

		// insert booking
		errInsert := insertBooking(user.Id.(int64), fromStr, "")
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}

		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	} else if !okFrom && okTo { // only clock out

		// check for valid dateTime Format
		if !checkDateTimeFormat(toStr) {
			log.Printf("ERROR: To string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: To string not in valid DateTime format.")
		}

		// insert "to" to existing booking
		log.Printf(`got "to" booking: %s`, toStr)

		// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if !res { // Error - user has no open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has no open booking")
		}

		errInsert := insertBooking(user.Id.(int64), "", toStr)
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}
		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	} else if okFrom && okTo { // full clocking

		// check for valid dateTime Format
		if !checkDateTimeFormat(fromStr) || !checkDateTimeFormat(toStr) {
			log.Printf("ERROR: From or To string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: From or To string not in valid DateTime format.")
		}

		// insert from and to
		log.Printf(`got full booking: %s - %s `, fromStr, toStr)

		// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if res { // Error - user has open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has already open booking")
		}

		errInsert := insertBooking(user.Id.(int64), fromStr, toStr)
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}
		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	}

	// error - got no values or other error
	return generateResponse("handleTimeBookingResponse", true, "Error: Missing / error timeBooking values.")

}

func handleGetBookings(content interface{}, user User) ([]byte, error) {
	bookings, err := dbGetBookings(user)
	if err != nil {
		return generateResponse("handleGetBookingsResponse", true, "Error: Error getting Bookings: "+err.Error())
	}

	return generateResponse("handleGetBookingsResponse", false, bookings)
}

func handleGetBookings_OLD(content interface{}) ([]byte, error) {
	contentStr, ok := content.(string)
	if !ok {
		return generateResponse("handleGetBookingsResponse", true, "Invalid content format for timebooking")
	}
	var response Message

	switch contentStr {
	case "currentMonth":
		response.Type = "getBookingsResponse"
		// Sample input string to be encoded
		input := `[{"from":"01.01.2020 00:00:00", "to":"01.01.2020 00:01:00"},{"from":"02.01.2020 00:00:00", "to":"02.01.2020 00:13:00"}]`

		// Define a slice of TimeRange to hold the parsed time ranges
		var timeRanges []TimeRange

		// Parse the input JSON string into the slice
		err := json.Unmarshal([]byte(input), &timeRanges)
		if err != nil {
			fmt.Println("Error parsing input:", err)
			//	return
		}

		// Create a Message object with the parsed time ranges as the content
		response = Message{
			Type:    "getBookingsResponse",
			Content: timeRanges,
		}
		return generateResponse("handleGetBookingsResponse", false, timeRanges)

	default:
		return generateResponse("handleGetBookingsResponse", true, "invalid getBookings message - use 'currentMonth' for example")

	}

	//return json.Marshal(response)
}
