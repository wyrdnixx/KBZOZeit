package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func processMessage_OLD(conn *websocket.Conn, msg Message) {

	log.Printf("Processing Message from Client %s : %s", conn.RemoteAddr(), msg)

	switch msg.Type {
	case "textMessage":
		log.Printf("Text message received: %s", msg.Content)
		err := conn.WriteMessage(websocket.TextMessage, []byte("{\"text\":\"got your message\"}"))
		if err != nil {
			log.Println("Write:", err)
			break
		}

	case "echoTest":
		log.Printf("Echo request message received: %s", msg.Content)
		responseJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println("Write:", err)
			break
		}
		err = conn.WriteMessage(websocket.TextMessage, responseJSON)
		if err != nil {
			log.Println("Write:", err)
			break
		}

	}
}

// processMessage handles the incoming message based on its "type"
func processMessage(msg []byte) ([]byte, error) {
	var message Message

	// Try to unmarshal the incoming message into the Message struct
	if err := json.Unmarshal(msg, &message); err != nil {
		return generateErrorJSON("Invalid JSON format on message")
	}

	// Switch on the "type" field to handle different types of messages
	switch message.Type {
	case "echoTest":
		// Handle the echoTest case where content is a string
		return handleEchoTest(message.Content)
		// Example case: Echo the content back
		//response := map[string]string{
		//	"type":    "echoResponse",
		//	"content": message.Content,
		//}
		//return json.Marshal(response)

	case "timebooking":
		// Handle the timebooking case where content is an object
		return handleTimeBooking(message.Content)

	case "clocking":
		return handleClocking(message.Content)

	// You can add more cases here for different message types
	default:
		// Unsupported type
		return generateErrorJSON(fmt.Sprintf("Unsupported message type: %s", message.Type))
	}
}

// generateErrorJSON creates an error JSON response
func generateErrorJSON(errorMessage string) ([]byte, error) {
	errorResponse := ErrorResponse{
		Type:    "error",
		Message: errorMessage,
	}
	return json.Marshal(errorResponse)
}
