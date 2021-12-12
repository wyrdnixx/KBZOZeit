package database

import (
	"database/sql"
	"errors"

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
		Name VARCHAR(40) not null,			
		Enabled INT,			
		primary key (Name)
	);`
	crDevices := `
	CREATE TABLE IF NOT EXISTS Devices (			
		Id VARCHAR(40) not null,
		FUsers VARCHAR(40) NOT NULL,		
		Enabled INT,			
		primary key (id)
	);`
	crTimeAccounting := `
	CREATE TABLE IF NOT EXISTS TimeAccounting (
		Uuid INT auto_increment,			
		FUsers VARCHAR(40) NOT NULL,
		FromDate  TIMESTAMP,
		ToDate   TIMESTAMP,			
		primary key (Uuid)
	);`

	QueryDB(crUsers)
	QueryDB(crTimeAccounting)
	QueryDB(crDevices)
	//	AddUser("EnabledExampleUser")
	//	DisableUser("EnabledExampleUser")
	//	FindUser("EnabledExampleUser")
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
func AddUser(name string) error {
	utils.Log(1, "AddUser()", ": "+name)
	db, err := sql.Open("mysql", cfg.DB_USERNAME+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)
	if err != nil {
		utils.Log(4, "AddUser()", "Could not connect to database: "+err.Error())
		return err
	} else {
		sqlSelect := "INSERT INTO Users (Name, Enabled) VALUES ( '" + name + "', 1)"
		utils.Log(1, "DisableUser()", "SQL-Query: "+sqlSelect)
		_, err := db.Query(sqlSelect)
		if err != nil {
			utils.Log(3, "AddUser()", "Error adding User: "+err.Error())
			return nil

		} else {
			utils.Log(1, "AddUser()", "User created... ")
			return nil
		}
	}
}

func DisableUser(name string) error {
	utils.Log(1, "DisableUser()", ": "+name)
	db, err := sql.Open("mysql", cfg.DB_USERNAME+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)
	if err != nil {
		utils.Log(4, "DisableUser()", "Could not connect to database: "+err.Error())
		return err
	} else {
		sqlSelect := "UPDATE Users Set Enabled = 0 where Name = '" + name + "'"
		utils.Log(1, "DisableUser()", "SQL-Query: "+sqlSelect)
		_, err := db.Query(sqlSelect)
		if err != nil {
			utils.Log(3, "DisableUser()", "Error disabling User: "+err.Error())
			return nil

		} else {
			utils.Log(1, "DisableUser()", "User has been disabled... ")
			return nil
		}
	}
}

func FindUser(name string) (models.User, error) {
	utils.Log(1, "FindUser()", ": "+name)
	var u models.User
	db, err := sql.Open("mysql", cfg.DB_USERNAME+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)
	if err != nil {
		utils.Log(4, "FindUser()", "Could not connect to database: "+err.Error())
		return u, err
	} else {
		sqlSelect := "Select Name, Enabled from Users where Name = '" + name + "'"
		utils.Log(1, "FindUser()", "SQL-Query: "+sqlSelect)

		err = db.QueryRow(sqlSelect).Scan(&u.Name, &u.Enabled)
		if err != nil {
			utils.Log(3, "FindUser()", "Error findung User: "+err.Error())
			return u, nil

		} else {
			utils.Log(1, "FindUser()", "User found... ")
			return u, nil
		}
	}
}

func GetUsers() (models.Users, error) {
	utils.Log(1, "GetUsers()", "Getting users from DB... ")

	users := models.Users{}

	db, err := sql.Open("mysql", cfg.DB_USERNAME+":"+cfg.DB_PASSWORD+"@tcp("+cfg.DB_HOST+":"+cfg.DB_PORT+")/"+cfg.DB_NAME)

	if err != nil {
		utils.Log(3, "GetUsers()", "Could not connect to database: "+err.Error())
		return users, err

	} else {

		rows, err := db.Query("select * from Users;")
		if err != nil {
			utils.Log(3, "GetUsers()", "Error getting users from DB: "+err.Error())
		} else {
			for rows.Next() {
				u := models.User{}
				err = rows.Scan(&u.Name, &u.Enabled)
				//fmt.Printf("\nGetFirmen -- Got: %s", itm)
				users.User = append(users.User, u)

			}
		}

	}
	defer db.Close()
	if len(users.User) == 0 {
		return users, errors.New(`No users got from DB`)
	}
	return users, nil
}
