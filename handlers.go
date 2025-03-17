package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
		// insert from to new booking
		log.Printf(`got "from" booking: %s`, fromStr)

		// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleClockingResponse", true, "Error: cannot check fo open bookings")
		}
		if res { // Error - user has open bookings
			return generateResponse("handleClockingResponse", true, "Error: User has already open booking")
		}

		// insert booking
		errInsert := insertBooking(user.Id.(int64), fromStr, "")
		if errInsert != nil {
			return generateResponse("handleClockingResponse", true, "Error: DB error: "+errInsert.Error())
		}

		return generateResponse("handleClockingResponse", false, "booking processed")

	} else if !okFrom && okTo { // only clock out
		// insert "to" to existing booking
		log.Printf(`got "to" booking: %s`, toStr)

		// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleClockingResponse", true, "Error: cannot check fo open bookings")
		}
		if !res { // Error - user has no open bookings
			return generateResponse("handleClockingResponse", true, "Error: User has no open booking")
		}

		errInsert := insertBooking(user.Id.(int64), "", toStr)
		if errInsert != nil {
			return generateResponse("handleClockingResponse", true, "Error: DB error: "+errInsert.Error())
		}
		return generateResponse("handleClockingResponse", false, "booking processed")

	} else if okFrom && okTo { // full clocking
		// insert from and to
		log.Printf(`got full booking: %s - %s `, fromStr, toStr)

		// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleClockingResponse", true, "Error: cannot check fo open bookings")
		}
		if res { // Error - user has open bookings
			return generateResponse("handleClockingResponse", true, "Error: User has already open booking")
		}

		// ToDo: insert booking
		return generateResponse("handleClockingResponse", false, "booking processed")
	}

	// error - got no values or other error
	return generateResponse("handleTimeBookingResponse", true, "Error: Missing / error timeBooking values.")

}

// handleTimeBooking processes timebooking type messages
func handleTimeBooking_OLD_TEST(content interface{}) ([]byte, error) {
	// First, convert content to a map
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return generateResponse("handleTimeBookingResponse", true, "Invalid content format for timebooking")
	}

	// Extract the "from" and "to" fields
	fromStr, okFrom := contentMap["from"].(string)
	toStr, okTo := contentMap["to"].(string)
	if !okFrom || !okTo {
		return generateResponse("handleTimeBookingResponse", true, "Missing 'from' or 'to' in timebooking content")
	}

	// Parse the datetime strings into time.Time objects
	layout := "02.01.2006 15:04:05"
	fromTime, err := time.Parse(layout, fromStr)
	if err != nil {
		return generateResponse("handleTimeBookingResponse", true, "Invalid 'from' datetime format")
	}
	toTime, err := time.Parse(layout, toStr)
	if err != nil {
		return generateResponse("handleTimeBookingResponse", true, "Invalid 'to' datetime format")
	}

	fmt.Printf("got booking from %s to %s\n", fromTime.String(), toTime.String())

	// Respond with a success message and the parsed time data
	response := map[string]interface{}{
		"type":    "timebookingResponse",
		"from":    fromTime.String(),
		"to":      toTime.String(),
		"message": "Time booking processed successfully",
	}
	return json.Marshal(response)
}

// OLD - do everything in handleTimeBooking
func handleClocking_OLD(content interface{}, user User) ([]byte, error) {
	contentStr, ok := content.(string)
	if !ok {
		return generateResponse("handleClockingResponse", true, "invalid string in contend")
	}
	var response Message

	log.Printf("clocking user: %s", user.Username)

	switch contentStr {
	case "clockIn":
		log.Println("clocking in")
		// ToDo DB ClockIn User
		//_, err := testInsert()

		res, err := getOpenBookings(user)

		if err != nil {
			return generateResponse("handleClockingResponse", true, "Error checking open bookings : "+err.Error())
			//response.Type = "clockingResponseError"
			//response.Content = "error clocking in processed : " + err.Error() + " - " + getcurrentTimestamp()
		} else if res { // if user has open booking
			return generateResponse("handleClockingResponse", true, "Error: User has already open booking")
		} else {
			// booking for user

			return generateResponse("handleClockingResponse", false, "booking successfull")
		}

	case "clockOut":
		log.Println("clocking out")
		response.Type = "clockingResponse"
		response.Content = "clocking out processed successfully" + getcurrentTimestamp()

	default:
		log.Println("invalid clocking message")
		response.Type = "clockingResponseError"
		response.Content = "invalid clocking message"

	}

	return json.Marshal(response)
}

func handleGetBookings(content interface{}) ([]byte, error) {
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
