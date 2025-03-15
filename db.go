package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(db *sql.DB) error {

	//defer db.Close()
	// Get the database file path
	filePath, err := getDatabaseFilePath(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database file path:", filePath)

	// Create a simple table
	//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT not null);`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "users" ("id" INTEGER NOT NULL,"name"	TEXT NOT NULL UNIQUE,"password" TEXT NOT NULL,"token" TEXT, "isClockedIn" INT, PRIMARY KEY("id"));`)
	if err != nil {
		log.Fatal("initDB: " + err.Error())
		return err
	}

	//inital insert default user

	insertTask := &DBTask{
		Action:   "insert",
		Query:    `INSERT INTO users (name,password, token,isClockedIn) VALUES (?,?,?,?);`,
		Args:     []interface{}{"admin", "admin", "adminToken", 0},
		Response: make(chan any),
	}

	result, err := dbEventBus.SubmitTask(insertTask)
	if err != nil {
		// Check for error if user already exists - "UNIQUE constraint failed"
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			log.Printf("inital Admin user already exists")
		} else {
			log.Printf("error creating inital Admin user: %s", err)
		}

	} else {
		log.Printf("inital Admin user created: %s", result)
	}

	return nil
}

// NewDBEventBus creates a new EventBus and starts the task worker.
func NewDBEventBus(db *sql.DB) *DBEventBus {
	log.Printf("creating eventbus\n")
	bus := &DBEventBus{
		db:    db,
		tasks: make(chan *DBTask, 100), // Buffered channel to hold tasks
	}
	bus.startWorker()
	return bus
}

// startWorker starts the task processor in a separate Goroutine.
func (bus *DBEventBus) startWorker() {
	go func() {
		for task := range bus.tasks {
			bus.processTask(task)
		}
	}()
}

/* // processTask processes a single task (either fetch or insert).
func (bus *DBEventBus) processTask(task *DBTask) {
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
} */

// processTask processes a single task (fetch or insert/update/delete).
func (bus *DBEventBus) processTask(task *DBTask) {
	switch task.Action {
	case "insert", "update", "delete":
		log.Printf("token: %s", task.Query)
		fmt.Printf("Executing Query: %s with Args: %v\n", task.Query, task.Args)
		result, err := bus.db.Exec(task.Query, task.Args...)
		if err != nil {
			task.Response <- err
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			task.Response <- err
			return
		}
		task.Response <- rowsAffected
	case "fetch":
		rows, err := bus.db.Query(task.Query, task.Args...)
		if err != nil {
			task.Response <- err
			return
		}
		task.Response <- rows

	case "fetchRow":
		// Handle select queries (fetch task)
		rows, err := bus.db.Query(task.Query, task.Args...)
		if err != nil {
			task.Response <- err
			return
		}
		task.Response <- rows
	}
}

// SubmitTask submits a task to the EventBus and waits for a result.
func (bus *DBEventBus) SubmitTask(task *DBTask) (any, error) {
	bus.tasks <- task
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
func (bus *DBEventBus) Close() {
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

func getUserbyToken(token string) (User, error) {
	// Fetch users
	fetchTask := &DBTask{
		Action:   "fetch",
		Query:    `SELECT id, name FROM users where token = (?);`,
		Args:     []interface{}{token},
		Response: make(chan any),
	}
	rowsResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return User{}, err
	}

	rows := rowsResult.(*sql.Rows)
	defer rows.Close()

	var user User
	fmt.Println("Fetched users:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	//fmt.Println("Fetched users from DB:", users)
	return user, nil
}

// Old method from Test using websocket
func dbCheckUserPasswd(username string, passwd string) (User, error) {
	// Fetch users
	fetchTask := &DBTask{
		Action:   "fetch",
		Query:    `SELECT id, name  FROM users where name = (?) and password = (?);`,
		Args:     []interface{}{username, passwd},
		Response: make(chan any),
	}
	rowsResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return User{}, err
	}

	rows := rowsResult.(*sql.Rows)
	defer rows.Close()

	var user User
	fmt.Println("Fetched users:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	//fmt.Println("Fetched users from DB:", users)
	return user, nil
}

// Function to validate user credentials from SQLite database
func validateUser(username, password string) bool {
	var storedPassword string
	// Fetch user
	fetchTask := &DBTask{
		Action:   "fetchRow",
		Query:    "SELECT password FROM users WHERE name = ?",
		Args:     []interface{}{username},
		Response: make(chan any),
	}
	rowResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return false
	}

	// Assert the response to *sql.Rows
	rows, ok := rowResult.(*sql.Rows)
	if !ok {
		log.Fatal("Failed to assert rowsResult to *sql.Rows")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&storedPassword); err != nil {
			log.Fatal(err)
		}
	}
	// ToDO: In a real-world app, use a hashed password comparison, not plaintext
	return password == storedPassword
}

func dbUpdateToken(username string, token string) error {
	// update user token
	//id := strconv.FormatInt(userID, 10)
	insertTask := &DBTask{
		Action: "update",
		Query:  `UPDATE users SET token = ? WHERE name = ?;`,
		Args:   []interface{}{token, username},
		//Args:     []interface{}{token},
		Response: make(chan any),
	}
	var rowsAffected int64
	result, err := dbEventBus.SubmitTask(insertTask)
	if err != nil {
		return err
	}
	if result, ok := result.(sql.Result); ok {
		fmt.Println("Result is of type sql.Result")
		// Now you can use sqlResult
		rowsAffected, err = result.RowsAffected()
		if err == nil {
			fmt.Printf("Rows affected: %d\n", rowsAffected)

		}

	}
	return nil
}

func testInsert() (int64, error) {
	// Insert a user
	insertTask := &DBTask{
		Action:   "insert",
		Query:    `INSERT INTO users (name,password) VALUES (?,?);`,
		Args:     []interface{}{"123", "test"},
		Response: make(chan any),
	}
	var rowsAffected int64
	result, err := dbEventBus.SubmitTask(insertTask)
	if err != nil {
		return 0, err
	}
	if result, ok := result.(sql.Result); ok {
		fmt.Println("Result is of type sql.Result")
		// Now you can use sqlResult
		rowsAffected, err = result.RowsAffected()
		if err == nil {
			fmt.Printf("Rows affected: %d\n", rowsAffected)

		}

	}
	return rowsAffected, nil
}
