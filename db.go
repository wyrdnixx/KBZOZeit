package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// Task represents a database operation, either a fetch or an insert.
type Task struct {
	Action   string        // "fetch" or "insert"
	Query    string        // SQL query
	Args     []interface{} // Arguments for the query
	Response chan any      // Response channel for the result or error
}

// EventBus manages the processing of tasks sequentially.
type EventBus struct {
	db     *sql.DB
	tasks  chan *Task
	wg     sync.WaitGroup
	closed bool
	mu     sync.Mutex
}

func initDB(db *sql.DB) error {

	//defer db.Close()
	// Get the database file path
	filePath, err := getDatabaseFilePath(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database file path:", filePath)

	// Create the EventBus
	//eventBus := NewEventBus(db)

	// Create a simple table

	//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT not null);`)
	/* _, err = db.Exec(`CREATE TABLE IF NOT EXISTS "users" ("id" INTEGER NOT NULL,"name"	TEXT NOT NULL,"password" TEXT NOT NULL,PRIMARY KEY("id"));`)
	if err != nil {
		log.Fatal("initDB: " + err.Error())
		return err
	} */
	return nil
}

// NewEventBus creates a new EventBus and starts the task worker.
func NewEventBus(db *sql.DB) *EventBus {
	log.Printf("creating eventbus\n")
	bus := &EventBus{
		db:    db,
		tasks: make(chan *Task, 100), // Buffered channel to hold tasks
	}
	bus.startWorker()
	return bus
}

// startWorker starts the task processor in a separate Goroutine.
func (bus *EventBus) startWorker() {
	go func() {
		for task := range bus.tasks {
			bus.processTask(task)
		}
	}()
}

// processTask processes a single task (either fetch or insert).
func (bus *EventBus) processTask(task *Task) {
	log.Printf("new processTask: %s : %s : %s \n", task.Action, task.Query, task.Args)

	defer bus.wg.Done()

	switch task.Action {
	case "insert":
		result, err := bus.db.Exec(task.Query, task.Args...)
		if err != nil {
			task.Response <- err
			return
		}
		task.Response <- result
	case "fetch":
		rows, err := bus.db.Query(task.Query, task.Args...)
		if err != nil {
			task.Response <- err
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		cols, err := rows.Columns()
		if err != nil {
			task.Response <- err
			return
		}
		for rows.Next() {
			rowData := make([]interface{}, len(cols))
			rowPointers := make([]interface{}, len(cols))
			for i := range rowData {
				rowPointers[i] = &rowData[i]
			}
			if err := rows.Scan(rowPointers...); err != nil {
				task.Response <- err
				return
			}
			rowMap := make(map[string]interface{})
			for i, col := range cols {
				rowMap[col] = rowData[i]
			}
			results = append(results, rowMap)
		}
		task.Response <- results
	}
}

// SubmitTask submits a task to the EventBus and waits for a result.
func (bus *EventBus) SubmitTask(task *Task) (any, error) {
	log.Printf("new db-task submitted\n")
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.closed {
		bus.mu.Unlock()
		return nil, fmt.Errorf("event bus is closed")
	}
	bus.wg.Add(1)

	bus.tasks <- task

	// bus.mu.Unlock()

	// Wait for the response
	response := <-task.Response
	switch v := response.(type) {
	case error:
		return nil, v
	default:
		return v, nil
	}
}

// Close closes the EventBus and waits for all tasks to complete.
func (bus *EventBus) Close() {
	bus.mu.Lock()
	if bus.closed {
		bus.mu.Unlock()
		return
	}
	bus.closed = true
	close(bus.tasks)
	bus.mu.Unlock()
	bus.wg.Wait()
}

func getDatabaseFilePath(db *sql.DB) (string, error) {
	var name, file string
	var seq int
	// Execute PRAGMA database_list to get the database filename
	row := db.QueryRow("PRAGMA database_list;")
	err := row.Scan(&seq, &name, &file)
	if err != nil {
		return "", err
	}
	return file, nil
}
