package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/wyrdnixx/KBZOZeit/models"
	"github.com/wyrdnixx/KBZOZeit/utils.go"
)

var cfg models.Configuration

func init() {
	cfg = utils.GetConfig()
}

func Initdb() {

	utils.Log(1, "InitDB()", "starting initiating database")

	db, err := sql.Open("mysql", cfg.DB_USERNAME+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)

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
		if errUsers != nil {
			utils.Log(4, "InitDB()", "Could not create table Users: "+errUsers.Error())

		} else {
			utils.Log(1, "InitDB()", "init table Users completed")
		}

		_, errTimeAccounting := db.Query(`
		CREATE TABLE IF NOT EXISTS TimeAccounting (
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
