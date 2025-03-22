package main

import (
	"encoding/json"
	"log"
	"math"
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

	if okFrom && !okTo { /////////////////////////////////////////// only clock in

		// check for valid dateTime Format
		if !checkDateTimeFormat(fromStr) {
			log.Printf("ERROR: From string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: From string not in valid DateTime format.")
		}
		// insert from to new booking
		log.Printf(`got "from" booking: %s`, fromStr)

		// Check if user has alrady open bookings
		//ToDO - Fehler - wird immer erstellt
		openBooking, err := getOpenBookings(user)
		var zeroBooking Booking
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if openBooking != zeroBooking { // Error - user has open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has already open booking")
		}

		// insert booking
		errInsert := insertBooking(user.Id.(int64), fromStr, "", "")
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}

		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	} else if !okFrom && okTo { ///////////////////////////////////////// only clock out

		// check for valid dateTime Format
		if !checkDateTimeFormat(toStr) {
			log.Printf("ERROR: To string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: To string not in valid DateTime format.")
		}

		// insert "to" to existing booking
		log.Printf(`got "to" booking: %s`, toStr)

		/* 	// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if !res { // Error - user has no open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has no open booking")
		} */

		// Check if user has alrady open bookings
		openBooking, err := getOpenBookings(user)
		var zeroBooking Booking
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}

		if openBooking == zeroBooking {
			log.Printf("Error: No open booking found. ")
			return generateResponse("handleTimeBookingResponse", true, "Error: User has no open booking")
		}

		duration, err := calcDuration(openBooking.From, toStr)
		if err != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: Error parsing dates for duration calculation")
		}
		// ToDo: calculate duration
		errInsert := insertBooking(user.Id.(int64), "", toStr, duration)
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}
		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	} else if okFrom && okTo { ///////////////////////////////////////// full clocking

		// check for valid dateTime Format
		if !checkDateTimeFormat(fromStr) || !checkDateTimeFormat(toStr) {
			log.Printf("ERROR: From or To string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: From or To string not in valid DateTime format.")
		}

		// insert from and to
		log.Printf(`got full booking: %s - %s `, fromStr, toStr)

		// Check if user has alrady open bookings
		openBooking, err := getOpenBookings(user)
		var zeroBooking Booking
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}

		if openBooking != zeroBooking {
			log.Printf("Error: found open booking from: %s", openBooking.From)
			return generateResponse("handleTimeBookingResponse", true, "Error: User has already open booking")
		}

		duration, err := calcDuration(fromStr, toStr)
		if err != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: Error parsing dates for duration calculation")
		}
		/* //calculate duration
		layout := "02.01.2006 15:04"
		startTime, err1 := time.Parse(layout, fromStr)
		endTime, err2 := time.Parse(layout, toStr)
		// Check for parsing errors
		if err1 != nil || err2 != nil {
			log.Printf("Error parsing dates: %s - %s", err1, err2)
			return generateResponse("handleTimeBookingResponse", true, "Error: Error parsing dates for duration calculation")
		}
		// Calculate the duration between the two times
		timeDuration := endTime.Sub(startTime)
		// Convert the duration to hours and minutes
		totalHours := timeDuration.Hours()
		hours := int(totalHours)                           // Get the integer part (hours)
		minutes := int((totalHours - float64(hours)) * 60) // Get the fractional part as minutes

		// Format the duration as "hours.minutes"
		duration := fmt.Sprintf("%d.%02d", hours, minutes) */

		errInsert := insertBooking(user.Id.(int64), fromStr, toStr, duration)
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

func handleStartRecalc(content interface{}, user User) ([]byte, error) {

	// ToDo: Fehler - summe aus SQL Select sind nicht minuten (60er Basis)
	var monthSollZeit float64 = 10
	var monthCount float64 = 5
	log.Printf("Sollzeit: %f", monthSollZeit)
	var sumSoll = monthSollZeit * monthCount
	log.Printf("Sollzeit gesamt: %f", sumSoll)
	sumIst, _ := getFullTimeAccountings(user)
	log.Printf("Istzeit gesamt: %f", sumIst)

	dif := sumIst - sumSoll
	valDif := math.Round(dif*100) / 100
	log.Printf("Dif: %f", valDif)

	return generateResponse("handleStartRecalcResponse", false, "recalc finished")
}
