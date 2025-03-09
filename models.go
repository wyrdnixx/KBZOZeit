package main

import "github.com/gorilla/websocket"

// Message represents the structure of a WebSocket message
type Message struct {
	Type string `json:"type"` // Type of the message (e.g., "text", "notification")
	//Content string `json:"content"` // Actual message content
	Content interface{} `json:"content"`
}

type timebookingMessage struct {
	From string `json:"from"` // Actual message content
	To   string `json:"to"`   // Actual message content
}

type WSConnections struct {
	C []*websocket.Conn
}

// ErrorResponse is the structure for sending error messages
type ErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Define the structure for the embedded content (time ranges)
type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}
