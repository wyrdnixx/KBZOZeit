package database

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"

	"github.com/wyrdnixx/KBZOZeit/utils.go"
)

func init() {
	cfg = utils.GetConfig()
}

func TESTQueryDB(query string) sql.Rows {
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

func AddDevice(username string, device string) error {
	utils.Log(1, "AddDevice()", ": "+device+" - to User -> "+username)

	u, err := FindUser(username)
	if err != nil {
		return err
	} else {
		if u.Enabled == 1 {
			utils.Log(1, "AddDevice()", ": user "+u.Name+" is valid, creating device...")
			return nil
		} else {
			return errors.New("user is disabled")
		}

	}
}
