package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "users" ("id" INTEGER NOT NULL, "name"	TEXT NOT NULL UNIQUE,"pwdHash" TEXT NOT NULL, "token" TEXT, PRIMARY KEY("id"));`)
	if err != nil {
		log.Fatal("initDB: " + err.Error())
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "bookings" ("id" INTEGER NOT NULL,"userId" INTEGER NOT NULL,"from" TEXT NOT NULL, "to" TEXT, "duration" TEXT, PRIMARY KEY("id"));`)
	if err != nil {
		log.Fatal("initDB: " + err.Error())
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "employee" ("id" INTEGER NOT NULL,"userId" INTEGER NOT NULL,"from" TEXT NOT NULL, "to" TEXT, "hoursPerMonth" TEXT, PRIMARY KEY("id"));`)
	if err != nil {
		log.Fatal("initDB: " + err.Error())
		return err
	}
	//inital insert default user

	_, err = db.Exec(`INSERT INTO "employee" ("userID","from", "to","hoursPerMonth") VALUES ("1","01.08.2023","","10");`)
	if err != nil {
		log.Fatal("initDB: " + err.Error())
		return err
	}

	insertTask := &DBTask{
		Action:   "insert",
		Query:    `INSERT INTO users (name, pwdHash, token,isClockedIn) VALUES (?,?,?,?);`,
		Args:     []interface{}{"admin", "testhash", "adminToken", 0},
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

// processTask processes a single task (fetch or insert/update/delete).
func (bus *DBEventBus) processTask(task *DBTask) {
	switch task.Action {
	case "insert", "update", "delete":
		//log.Printf("token: %s", task.Query)
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
	//log.Printf("Query finished")
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

	for rows.Next() {
		//var id int
		//var name string
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			log.Fatal("User validation using token error: " + err.Error())

		}
		log.Printf("User validated - ID: %d, Name: %s\n", user.Id, user.Username)
	}
	//fmt.Println("Fetched users from DB:", users)
	return user, nil
}

// func getOpenBookings(user User) (bool, error) {
func getOpenBookings(user User) (Booking, error) {

	fetchTask := &DBTask{
		Action:   "fetch",
		Query:    `SELECT id, "from", "to", "duration" FROM bookings WHERE userId = (?) AND "to" IS NULL ;`,
		Args:     []interface{}{user.Id},
		Response: make(chan any),
	}
	rowsResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return Booking{}, err
	}

	rows := rowsResult.(*sql.Rows)
	defer rows.Close()

	if rows.Next() {
		log.Printf("Error: found open booking...")
		var booking Booking // Create a variable for each row
		if err := rows.Scan(&booking.Id, &booking.From, &booking.To, &booking.Duration); err != nil {
			log.Printf("Error get bookings for user : %s : %s", user, err)
		}
		return booking, nil
	} else {
		return Booking{}, nil
	}
}

func dbGetBookings(user User) ([]Booking, error) {

	fetchTask := &DBTask{
		Action:   "fetch",
		Query:    `SELECT id, "from", "to" FROM bookings WHERE userId = (?) ;`,
		Args:     []interface{}{user.Id},
		Response: make(chan any),
	}
	rowsResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	rows := rowsResult.(*sql.Rows)
	defer rows.Close()

	var allBookings []Booking

	for rows.Next() {
		var booking Booking // Create a variable for each row
		if err := rows.Scan(&booking.Id, &booking.From, &booking.To); err != nil {
			log.Printf("Error get bookings for user : %s : %s", user, err)
		}

		//log.Printf("UserID: %s Got booking from:  %s to: %s\n", user.Id, booking.From, booking.To)
		allBookings = append(allBookings, booking)
	}
	return allBookings, nil
}

func BookingIn(user User, from string) error {

	insertTask := &DBTask{
		Action: "insert",
		Query: `INSERT INTO "main"."bookings" ("userId", "from") VALUES (?, ?);
`,
		Args:     []interface{}{user.Id, from},
		Response: make(chan any),
	}

	var rowsAffected int64
	result, err := dbEventBus.SubmitTask(insertTask)
	log.Printf("Finshed")
	if err != nil {
		return err
	}
	if result, ok := result.(sql.Result); ok {
		rowsAffected, err = result.RowsAffected()
		if err == nil {
			log.Printf("Booking in DB ok - Rows affected: %d\n", rowsAffected)
			return nil
		}
	}
	return err

}

// Function to validate user credentials from SQLite database
func dbValidateUser(username, pwdHash string) (User, error) {

	var loginUser User
	// Fetch user
	fetchTask := &DBTask{
		Action:   "fetchRow",
		Query:    "SELECT id, name FROM users WHERE name = ? and pwdHash = ?",
		Args:     []interface{}{username, pwdHash},
		Response: make(chan any),
	}
	rowResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return User{}, err
	}

	// Assert the response to *sql.Rows
	rows, ok := rowResult.(*sql.Rows)
	if !ok {
		log.Fatal("Failed to assert rowsResult to *sql.Rows")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&loginUser.Id, &loginUser.Username); err != nil {
			log.Fatal(err)
		}
	}
	// ToDO: In a real-world app, use a hashed password comparison, not plaintext
	return loginUser, nil
}

// OLD - Test sing websocket json
func validateUser(username, pwdHash string) bool {
	var storedpwdHash string
	// Fetch user
	fetchTask := &DBTask{
		Action:   "fetchRow",
		Query:    "SELECT pwdHash FROM users WHERE name = ?",
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
		if err := rows.Scan(&storedpwdHash); err != nil {
			log.Fatal(err)
		}
	}
	// ToDO: In a real-world app, use a hashed password comparison, not plaintext
	return pwdHash == storedpwdHash
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
		//fmt.Println("Result is of type sql.Result")
		// Now you can use sqlResult
		rowsAffected, err = result.RowsAffected()
		if err == nil {
			fmt.Printf("Rows affected: %d\n", rowsAffected)
		}

	}
	return rowsAffected, nil
}

func insertBooking(userId int64, from string, to string, duration string) error {

	if from != "" && to == "" { // only clocking in

		insertTask := &DBTask{
			Action:   "insert",
			Query:    `INSERT INTO bookings ("userID","from") VALUES (?,?);`,
			Args:     []interface{}{userId, from},
			Response: make(chan any),
		}

		_, err := dbEventBus.SubmitTask(insertTask)
		if err != nil {
			log.Printf("Error executing insert booking: %s", err)
			return err
		} else {
			return nil
		}
	} else if from == "" && to != "" { // only clocking out

		insertTask := &DBTask{
			Action:   "update",
			Query:    `UPDATE bookings set "to" = ?, "duration" = ?  where userId = ? and "to" is null ;`,
			Args:     []interface{}{to, duration, userId, to},
			Response: make(chan any),
		}

		_, err := dbEventBus.SubmitTask(insertTask)
		if err != nil {
			log.Printf("Error executing update booking: %s", err)
			return err
		} else {
			return nil
		}

	} else if from != "" && to != "" { // full timeBooking
		insertTask := &DBTask{
			Action:   "insert",
			Query:    `INSERT INTO bookings ("userID","from", "to", "duration") VALUES (?,?,?,?);`,
			Args:     []interface{}{userId, from, to, duration},
			Response: make(chan any),
		}

		_, err := dbEventBus.SubmitTask(insertTask)
		if err != nil {
			log.Printf("Error executing insert booking: %s", err)
			return err
		} else {
			return nil
		}
	}

	return fmt.Errorf("ERROR: unexpected error time clocking")

}

func getFullTimeAccountings(user User) (float64, error) {

	fetchTask := &DBTask{
		Action:   "fetch",
		Query:    ` select "from","to" from bookings WHERE userId = (?);`,
		Args:     []interface{}{user.Id},
		Response: make(chan any),
	}
	rowsResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return 0, err
	}

	rows := rowsResult.(*sql.Rows)
	defer rows.Close()
	// Time layout based on your datetime format
	const layout = "02.01.2006 15:04"

	// Variables to accumulate total hours and minutes
	var totalMinutes int64

	for rows.Next() {
		var fromStr, toStr string

		// Scan the "from" and "to" values from the query
		if err := rows.Scan(&fromStr, &toStr); err != nil {
			return 0, err
		}

		// Parse the "from" time string
		fromTime, err := time.Parse(layout, fromStr)
		if err != nil {
			return 0, err
		}

		// Parse the "to" time string
		toTime, err := time.Parse(layout, toStr)
		if err != nil {
			return 0, err
		}

		// Calculate the duration between "from" and "to"
		duration := toTime.Sub(fromTime)

		// Add the duration to totalMinutes (convert duration to minutes)
		totalMinutes += int64(duration.Minutes())
	}

	// Check for any error that occurred during iteration
	if err := rows.Err(); err != nil {
		return 0, err
	}

	// Convert totalMinutes to hours and minutes
	totalHours := totalMinutes / 60
	remainingMinutes := totalMinutes % 60

	// Calculate the total as hours.minutes
	hoursMinutes := float64(totalHours) + float64(remainingMinutes)/60
	log.Printf("getFullTimeAccountings calculated sum: %f", hoursMinutes)

	return hoursMinutes, nil
}

func getEmployeementMonths(user User) (float64, float64, error) {

	fetchTask := &DBTask{
		Action: "fetch",
		Query: `SELECT 
    				hoursPerMonth,
        			(strftime('%Y', date('now')) - substr("from", 7, 4)) * 12 + (strftime('%m', date('now')) - substr("from", 4, 2)) AS months_passed
					FROM employee
					WHERE "from" <= date('now')
					AND userId = (?);`,
		Args:     []interface{}{user.Id},
		Response: make(chan any),
	}
	rowsResult, err := dbEventBus.SubmitTask(fetchTask)
	if err != nil {
		//log.Fatal(err)
		return 0, 0, err
	}

	rows := rowsResult.(*sql.Rows)
	defer rows.Close()

	var hoursPM float64
	var months_passed float64

	// Iterate over the rows and append the results to the slice
	for rows.Next() {
		err := rows.Scan(&hoursPM, &months_passed)
		if err != nil {
			log.Printf("Error getting getEmployeementMonths: %s ", err)
			return 0, 0, err
		}
	}
	//fmt.Println("DB got full time accountings of:", sum)

	return hoursPM, months_passed, nil
}
