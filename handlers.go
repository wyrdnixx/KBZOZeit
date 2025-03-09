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
		return generateErrorJSON("Invalid content for echoTest")
	}

	// Respond with the same content
	response := map[string]string{
		"type":    "echoResponse",
		"content": contentStr,
	}
	return json.Marshal(response)
}

// handleTimeBooking processes timebooking type messages
func handleTimeBooking(content interface{}) ([]byte, error) {
	// First, convert content to a map
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return generateErrorJSON("Invalid content format for timebooking")
	}

	// Extract the "from" and "to" fields
	fromStr, okFrom := contentMap["from"].(string)
	toStr, okTo := contentMap["to"].(string)
	if !okFrom || !okTo {
		return generateErrorJSON("Missing 'from' or 'to' in timebooking content")
	}

	// Parse the datetime strings into time.Time objects
	layout := "02.01.2006 15:04:05"
	fromTime, err := time.Parse(layout, fromStr)
	if err != nil {
		return generateErrorJSON("Invalid 'from' datetime format")
	}
	toTime, err := time.Parse(layout, toStr)
	if err != nil {
		return generateErrorJSON("Invalid 'to' datetime format")
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
		return generateErrorJSON("Invalid content for echoTest")
	}
	var response Message
	switch contentStr {
	case "clockIn":
		log.Println("clocking in")

		response.Type = "clockingResponse"
		response.Content = "clocking in processed successfully" + getcurrentTimestamp()

	case "clockOut":
		log.Println("clocking out")
		response.Type = "clockingResponse"
		response.Content = "clocking out processed successfully" + getcurrentTimestamp()

	default:
		log.Println("invalid clocking message")
		response.Type = "clockingResponseError"
		response.Content = "invalid clocking message" + getcurrentTimestamp()

	}

	return json.Marshal(response)
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
