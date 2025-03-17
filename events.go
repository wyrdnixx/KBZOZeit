package main

import (
	"encoding/json"
	"log"
)

// processMessage handles the incoming message based on its "type"
func processMessage(msg []byte, user User) ([]byte, error) {
	var message Message

	// Try to unmarshal the incoming message into the Message struct
	if err := json.Unmarshal(msg, &message); err != nil {
		return generateResponse("processMessageResponse", true, "Invalid JSON format on message")
	}

	log.Printf("websocket message from user: %s : %s", user, message)

	// Switch on the "type" field to handle different types of messages
	switch message.Type {
	case "echoTest":
		// Handle the echoTest case where content is a string
		return handleEchoTest(message.Content)

	// handleLogin jetzt aus websocket raus genommen
	/* case "login":
	return handleLogin(message.Content) */

	case "timebooking":
		// Handle the timebooking case where content is an object
		return handleTimeBooking(message.Content, user)

	/* case "clocking":
	return handleClocking(message.Content, user) */

	case "getBookings":
		return handleGetBookings(message.Content)

	// You can add more cases here for different message types
	default:
		// Unsupported type
		//return generateResponse(fmt.Sprintf("Unsupported message type: %s", message.Type))
		return generateResponse("processMessageResponse", true, "Unsupported message type: "+message.Type)
	}
}

// generateErrorJSON creates an error JSON response
func generateResponse(MsgType string, isError bool, responseMessage interface{}) ([]byte, error) {
	Response := Response{
		Type:      MsgType,
		IsError:   isError,
		Timestamp: getcurrentTimestamp(),
		Message:   responseMessage,
	}
	return json.Marshal(Response)
}
