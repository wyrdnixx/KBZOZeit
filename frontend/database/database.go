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
	utils.Log(1, "InitDB()", "starting initiating database....")
	crUsers := `
	CREATE TABLE IF NOT EXISTS Users (			
		Id INT auto_increment,		
		Name VARCHAR(40) not null,			
		Enabled INT,			
		primary key (id)
	);`
	crDevices := `
	CREATE TABLE IF NOT EXISTS Devices (			
		Id VARCHAR(40) not null,
		FUsers INT NOT NULL,		
		Enabled INT,			
		primary key (id)
	);`
	crTimeAccounting := `
	CREATE TABLE IF NOT EXISTS TimeAccounting (
		Id INT auto_increment,			
		FUsers INT NOT NULL,
		FromDate  TIMESTAMP,
		ToDate   TIMESTAMP,			
		primary key (Id)
	);`

	QueryDB(crUsers)
	QueryDB(crTimeAccounting)
	QueryDB(crDevices)
}

func QueryDB(query string) sql.Rows {
	utils.Log(1, "QueryDB()", "starting query database: "+query)

	db, err := sql.Open("mysql", cfg.DB_USERNAME+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)

	if err != nil {
		//log.Fatal(4, "ERROR: DB Connection: "+err.Error())
		utils.Log(4, "QueryDB()", "Could not connect to database: "+err.Error())
		return sql.Rows{}
	} else {
		//res, err := db.Query("select * from bla;")

		res, errQuery := db.Query(query)
		if errQuery != nil {
			utils.Log(4, "QueryDB()", "Could not query database: "+errQuery.Error())
		} else {
			utils.Log(1, "QueryDB()", "query finished")
			return *res
		}

	}
	defer db.Close()
	return sql.Rows{}
}
