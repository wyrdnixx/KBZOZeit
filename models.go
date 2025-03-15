package main

import (
	"database/sql"
	"sync"

	"github.com/gorilla/websocket"
)

// Message represents the structure of a WebSocket message
type Message struct {
	Type string `json:"type"` // Type of the message (e.g., "text", "notification")
	//User string `json:"user"` // wird eigentlich nicht benötigt - identifikation über token
	//Content string `json:"content"` // Actual message content
	Content interface{} `json:"content"`
}

type timebookingMessage struct {
	From string `json:"from"` // Actual message content
	To   string `json:"to"`   // Actual message content
}

type User struct {
	Id       interface{} `json:"id"`   // UserID
	Username interface{} `json:"name"` // Username
}

type LoginUser struct {
	Username string `json:"username"`
	PwdHash  string `json:"pwdHash"`
}

type WSConnections struct {
	C []*websocket.Conn
}

// ErrorResponse is the structure for sending error messages
type Response struct {
	Type      string      `json:"type"`
	IsError   bool        `json:"isError"`
	Timestamp string      `json:"timestamp"`
	Message   interface{} `json:"message"`
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

// DBTask represents a database operation, either a fetch or an insert.
type DBTask struct {
	Action   string        // "fetch" or "insert"
	Query    string        // SQL query
	Args     []interface{} // Arguments for the query
	Response chan any      // Response channel for the result or error
}

// DBEventBus manages the processing of tasks sequentially.
type DBEventBus struct {
	db     *sql.DB
	tasks  chan *DBTask
	wg     sync.WaitGroup
	closed bool
	mu     sync.Mutex
}
