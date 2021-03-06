package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/wyrdnixx/KBZOZeit/database"
	"github.com/wyrdnixx/KBZOZeit/models"
	"github.com/wyrdnixx/KBZOZeit/utils.go"
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	(*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	(*w).Header().Set("Content-Type", "application/json")

}

func TestApi(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	// w.Header().Set("Content-Type", "application/json")

	utils.Log(1, "TestApi() ", "got called: ")

	reqBody, err := ioutil.ReadAll(r.Body)
	utils.Log(1, "TestApi() ", "got called: "+string(reqBody))

	if err != nil {
		fmt.Printf("TestApi error: %s \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {

		m := models.ReceivedMessage{}
		errunmarsh := json.Unmarshal(reqBody, &m)
		if errunmarsh != nil && errunmarsh != io.EOF {
			utils.Log(3, "TestApi() ", "first unmarshal: "+string(reqBody)+errunmarsh.Error())
		} else {
			utils.Log(1, "TestApi() ", "got MsgType: "+m.MsgType)

			switch m.MsgType {
			case "test":
				utils.Log(1, "TestApi() ", "got messagetype Test")
			case "GetUsers":
				utils.Log(1, "TestApi() ", "got messagetype GetUsers")
				AdminGetUsers(w, r)
			case "LoginRequest":
				utils.Log(1, "TestApi() ", "got messagetype LoginRequest")
				LoginRequest(w, r, string(reqBody))
			case "TimeAccounting":
				utils.Log(1, "TestApi() ", "got messagetype TimeAccounting")
				TimeAccounting(w, r, string(reqBody))
			case "AddUserRequest":
				utils.Log(1, "TestApi() ", "got messagetype AddUserRequest")
				u := models.User{}
				err := json.Unmarshal(reqBody, &u)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"Result":"error json.unmarshal ` + err.Error() + `"}`))
				} else {
					AdminAddUser(w, r, u)
				}
			case "GetOpenTimeaccounting":
				utils.Log(1, "TestApi() ", "got messagetype GetOpenTimeaccounting")
				GetOpenTimeaccounting(w, r, string(reqBody))
			case "GetAccountings":
				utils.Log(1, "TestApi", "User requests Accountings: "+string(reqBody))
				res, _ := GetAccountings(w, r, string(reqBody))
				w.Write([]byte(res))
			default:
				utils.Log(3, "TestApi() ", "got unknown messagetype: "+m.MsgType)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"Result":"Got unknown message type ` + m.MsgType + `"}`))
			}

		}

	}
}

type AccountingEntry struct {
	Id       string  `json:"Id"`
	User     string  `json:"FUser"`
	FromDate *string `json:"FromDate"`
	ToDate   *string `json:"ToDate"`
}

func GetAccountings(w http.ResponseWriter, r *http.Request, msg string) (string, error) {
	m := models.GetAccountingsMessage{}

	err := json.Unmarshal([]byte(msg), &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Result":"error unmarshal GetAccountingsMessage ` + err.Error() + `"}`))
		return err.Error(), nil
	} else {
		utils.Log(2, "GetAccountings() ", " GetAccountings for user: "+m.User)
		s := "" //initialice emtpy
		if m.FromDate == "" && m.ToDate == "" {
			utils.Log(1, "GetAccountings() ", " GetAccountings all for user "+m.User)
			s = `select * from TimeAccounting where FUsers = "` + m.User + `";`
		}
		if m.FromDate != "" && m.ToDate == "" {
			utils.Log(1, "GetAccountings() ", " GetAccountings range only from: "+m.FromDate)
			s = `select * from TimeAccounting where FUsers = "` + m.User + `" and FromDate >= "` + m.FromDate + `";`
		}
		if m.FromDate == "" && m.ToDate != "" {
			utils.Log(1, "GetAccountings() ", " GetAccountings range only to: "+m.FromDate)
			s = `select * from TimeAccounting where FUsers = "` + m.User + `" and ToDate <= "` + m.ToDate + ` 23:59:59";`
		}
		if m.FromDate != "" && m.ToDate != "" {
			utils.Log(1, "GetAccountings() ", " GetAccountings range : "+m.FromDate+" to: "+m.ToDate)
			s = `select * from TimeAccounting where FUsers = "` + m.User + `" and FromDate >= "` + m.FromDate + `" and ToDate <= "` + m.ToDate + ` 23:59:59";`
		}

		fmt.Println("Select: " + s)
		rows := database.QueryDB(s)

		//js, err := json.Marshal(rows)
		if err != nil {
			utils.Log(2, "GetAccountings: Error marshal sql results: ", err.Error())
			return "", err
		} else {
			//utils.Log(2, "GetAccountings: SQL-Result: ", string(js))
			//	return string(js), nil
			var entrys = []AccountingEntry{}
			for rows.Next() {
				var e AccountingEntry
				if err := rows.Scan(&e.Id, &e.User, &e.FromDate, &e.ToDate); err != nil {
					return "", err
				}
				entrys = append(entrys, e)
			}

			js, err := json.Marshal(entrys)
			if err != nil {
				utils.Log(2, "GetAccountings", "error marshal sql results: "+err.Error())
				return "", err
			} else {
				utils.Log(1, "GetAccountings", "SQL-Result: "+string(js))
				return string(js), nil

			}
		}

	}
	//return "nil", nil
}

func GetOpenTimeaccounting(w http.ResponseWriter, r *http.Request, msg string) {

	// ToDo : reqest open accounting for user
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"Error":"GetOpenTimeaccounting error ` + msg + `"}`))
}

// OLD: find correct user and do the response
func RegisterIdentOld(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)

	reqBody, err := ioutil.ReadAll(r.Body)
	utils.Log(1, "RegisterIdent() ", " reqBody got : "+string(reqBody))
	//err = errors.New("holla die waldfee")
	if err != nil {

		utils.Log(1, "RegisterIdent() ", " error reading body: "+err.Error())
		str := `{"Error":"` + err.Error() + `"}`
		w.Write([]byte([]byte(str)))

	} else {

		m := models.UserDevice{}
		err := json.Unmarshal(reqBody, &m)
		if err != nil {
			utils.Log(1, "RegisterIdent() ", " error JSON Unmarshal: "+err.Error())
		} else {
			utils.Log(1, "RegisterIdent() ", " name to find: "+m.Name)
			res, err := database.FindUser(m.Name)
			if err != nil {
				utils.Log(1, "RegisterIdent() ", " not found in DB: "+res.Name)
			} else {
				utils.Log(1, "RegisterIdent() ", "found in DB: "+res.Name)
			}

		}

		w.Write(reqBody)
	}

}

func TimeAccounting(w http.ResponseWriter, r *http.Request, request string) {
	EnableCors(&w)

	m := models.TimeAccountingMessage{}
	err := json.Unmarshal([]byte(request), &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Result":"error unmarshal TimeAccountingMessage ` + err.Error() + `"}`))
	} else {
		switch m.Typ {
		case "Einstempeln":
			utils.Log(1, "TimeAccounting", "startAccounting for User: "+m.Name)
			err := database.StartTimeAccounting(m.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"Result":"error unmarshal TimeAccountingMessage ` + err.Error() + `"}`))
			} else {
				w.Write([]byte(`{"Result":"startAccounting successfully"}`))

			}
		case "getOpenTimer":
			utils.Log(1, "TimeAccounting", "User request open Timer: "+m.Name)
			x, _ := database.CheckOpenCounters(m.Name)
			bytes, _ := json.Marshal(x)
			utils.Log(1, "TimeAccounting", "found open timer:"+string(bytes))
			w.Write(bytes)
		case "Ausstempeln":
			utils.Log(1, "TimeAccounting", "stopAccounting for User: "+m.Name)
			rows, err := database.Ausstempeln(m)
			if err != nil || rows == 0 {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"Result":"error stopAccounting ` + err.Error() + `"}`))
			} else {
				//r :=strconv.FormatInt(rows,2)
				r := strconv.FormatInt(rows, 10)
				w.Write([]byte(`{"Result":"stopAccounting successfully, updated rows: ` + r + `"}`))
			}

		}
	}
}

func LoginRequest(w http.ResponseWriter, r *http.Request, request string) {
	EnableCors(&w)

	u := models.User{}
	err := json.Unmarshal([]byte(request), &u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Result":"error unmarshal loginRequest ` + err.Error() + `"}`))
	} else {
		res, err := database.FindUser(u.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"Result":"error on loginRequest ` + err.Error() + `"}`))
		} else if res.Enabled == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"Result":"user not existing or disabled"}`))
		} else {
			w.Write([]byte(`{"Result":"login successfully"}`))
		}
	}

}

func AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	u, err := database.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
	} else {
		response, err := json.Marshal(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.Write(response)
		}
	}

}

func AdminAddUser(w http.ResponseWriter, r *http.Request, u models.User) {
	EnableCors(&w)
	if u.Name == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.New(`{"Result":"no username set"}`).Error()))
	} else {
		err := database.AddUser(u.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			response := `{"Result":"user created"}`
			w.Write([]byte(response))
		}
	}

}
