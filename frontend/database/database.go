package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wyrdnixx/KBZOZeit/utils.go"
)

var DBUser string = "dbuser"
var DBPassword string = "dbpw"
var DBHost string = "docker"
var DBPort string = "3306"
var DBName string = "testdb"

func Initdb() {

	utils.Log(1, "InitDB()", "starting initiating database")

	db, err := sql.Open("mysql", DBUser+":"+DBPassword+"@tcp("+DBHost+":"+DBPort+")/"+DBName)

	if err != nil {
		//log.Fatal(4, "ERROR: DB Connection: "+err.Error())
		utils.Log(4, "InitDB()", "Could not connect to database: "+err.Error())
	} else {
		//res, err := db.Query("select * from bla;")

		_, errUsers := db.Query(`
		CREATE TABLE IF NOT EXISTS Users (			
			Id INT auto_increment,
			Uuid VARCHAR(40) not null,
			Name VARCHAR(40) not null,			
			Enabled INT,			
			primary key (id)
		);
		`)
		if err != nil {
			utils.Log(4, "InitDB()", "Could not create table Users: "+errUsers.Error())

		} else {
			utils.Log(1, "InitDB()", "init table Users completed")
		}

		_, errTimeAccounting := db.Query(`
		CREATE TABLE IF NOT EXIS2TS TimeAccounting (
			Id INT auto_increment,			
            FUsers INT NOT NULL,
			FromDate  TIMESTAMP,
			ToDate   TIMESTAMP,			
			primary key (Id)
		);
		`)
		if errTimeAccounting != nil {
			utils.Log(4, "InitDB()", "Could not create table TimeAccounting: "+errTimeAccounting.Error())

		} else {
			utils.Log(1, "InitDB()", "init table TimeAccounting completed")
		}

	}
	defer db.Close()
}
