package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
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

func handleLogin(content interface{}) ([]byte, error) {

	// Check contend of the authentication request
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return generateResponse("handleLoginResponse", true, "Invalid content format for timebooking")
	}
	log.Printf("handleLogin")
	// Extract the "from" and "to" fields
	username, okUsr := contentMap["username"].(string)
	passwrd, okPwd := contentMap["passwd"].(string)
	if !okUsr || !okPwd {
		return generateResponse("handleLoginResponse", true, "Problem with password from authenticatoin request")
	}

	// check user against database
	res, err := dbCheckUserPasswd(username, passwrd)
	log.Printf("back from dbCheckUserPasswd: %s", res)
	if err != nil {
		log.Printf("error checking user in DB: %s\n", err)
		return generateResponse("handleLoginResponse", true, "error checking user in DB")
	} else if res.Id == nil {
		return generateResponse("handleLoginResponse", true, "wrong username or password")
	} else {
		log.Printf("got User from DB: %s\n", res)

		// generate bearer Token and write to DB
		token, err := GenerateJWT()
		if err != nil {
			fmt.Println("Error generating token:", err)
			return generateResponse("handleLoginResponse", true, "Error generating token")
		} else {
			log.Printf("new token generated: %s\n", token)

			// TODO: update geht noch nicht - 0 rows affected
			_, err := dbUpdateToken(res.Id.(int64), token)
			if err != nil {
				log.Printf("error updating token: %s\n", err)
				return generateResponse("handleLoginResponse", true, "Error updating token ind DB: "+err.Error())
			} else {
				return generateResponse("handleLoginResponse", false, token)
			}

		}

		// return successfull authentication

	}

}

//ToDo
/* func dbUpdateUserToken(token string) err {
	updateTask := &DBTask{
		Action:   "update",
		Query:    `update INTO users (name,password, token,isClockedIn) VALUES (?,?,?,?);`,
		Args:     []interface{}{"admin", "admin", "adminToken", 0},
		Response: make(chan any),
	}

	result, err := dbEventBus.SubmitTask(insertTask)
} */

// handleTimeBooking processes timebooking type messages
func handleTimeBooking(content interface{}) ([]byte, error) {
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

func handleClocking(content interface{}) ([]byte, error) {
	contentStr, ok := content.(string)
	if !ok {
		return generateResponse("handleClockingResponse", true, "invalid string in contend")
	}
	var response Message
	switch contentStr {
	case "clockIn":
		log.Println("clocking in")
		_, err := testInsert()
		if err != nil {
			return generateResponse("handleClockingResponse", true, "error clocking in processed : "+err.Error())
			//response.Type = "clockingResponseError"
			//response.Content = "error clocking in processed : " + err.Error() + " - " + getcurrentTimestamp()
		} else {
			return generateResponse("handleClockingResponse", false, "clocking in processed successfully ")
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

func getcurrentTimestamp() string {
	// Define the layout
	layout := "02.01.2006 15:04:05"

	// Get the current date and time
	currentTime := time.Now()

	// Format the current time using the specified layout
	formattedTime := currentTime.Format(layout)
	return formattedTime
}
