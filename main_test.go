package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Example struct that will be sent and received as JSON
/* type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
} */

// Test function to send and receive a JSON WebSocket message
func TestWebSocket(t *testing.T) {
	// Setup the HTTP server and WebSocket handler
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebSocket)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Connect to the WebSocket server
	wsURL := "ws" + server.URL[4:] + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	// Create a test message
	msg := Message{
		Type:    "echoTest",
		Content: "Hello, WebSocket!",
	}

	// Send the message as JSON over the WebSocket connection
	err = conn.WriteJSON(msg)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Receive the response
	var receivedMsg Message
	err = conn.ReadJSON(&receivedMsg)
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	// Assert that the received message matches the sent message
	assert.Equal(t, msg, receivedMsg, "The received message should match the sent message.")
}
