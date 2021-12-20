package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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
			utils.Log(3, "TestApi() ", "unmarshal: "+errunmarsh.Error())
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
			default:
				utils.Log(3, "TestApi() ", "got unknown messagetype: "+m.MsgType)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"Result":"Got unknown message type ` + m.MsgType + `"}`))
			}

		}

	}
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
	type TimeAccountingMessage struct {
		MsgType string `json:MsgType`
		Name    string `json:Name`
		Typ     string `json:Typ`
	}

	m := TimeAccountingMessage{}
	err := json.Unmarshal([]byte(request), &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Result":"error unmarshal TimeAccountingMessage ` + err.Error() + `"}`))
	} else {
		switch m.Typ {
		case "startAccounting":
			utils.Log(1, "TimeAccounting", "startAccounting for User: "+m.Name)
			err := database.StartTimeAccounting(m.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"Result":"error unmarshal TimeAccountingMessage ` + err.Error() + `"}`))
			} else {
				w.Write([]byte(`{"Result":"startAccounting successfully"}`))

			}
		case "stopAccounting":
			utils.Log(1, "TimeAccounting", "stopAccounting for User: "+m.Name)
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
