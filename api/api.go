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
			case "RegisterRequest":
				utils.Log(1, "TestApi() ", "got messagetype RegisterRequest")
				mRegisterReq := models.MsgRegisterRequest{}
				errRegister := json.Unmarshal(reqBody, &mRegisterReq)
				if errRegister != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"Result":"error json.unmarshal ` + errRegister.Error() + `"}`))
				} else {
					utils.Log(1, "TestApi() ", "RegisterRequest got Name: "+mRegisterReq.Name)
					utils.Log(1, "TestApi() ", "RegisterRequest got Uuid: "+mRegisterReq.Uuid)

					// ToDo - check for result
					err := RegisterIdent(mRegisterReq)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(`{"Result":"` + err.Error() + `"}`))
					} else {
						w.Write([]byte(`{"Result":"Processed"}`))
					}

				}
			default:
				utils.Log(3, "TestApi() ", "got unknown messagetype: "+m.MsgType)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"Result":"Got unknown message type ` + m.MsgType + `"}`))
			}

		}

	}
}

func RegisterIdent(m models.MsgRegisterRequest) error {
	// ToDo - error handling
	u, err := database.FindUser(m.Name)
	if (models.User{} == u) { //is empty
		// ok - user not found in DB
		database.AddUser(m.Name)
		database.AddDevice(m.Name, m.Uuid)
		return nil
	} else if (models.User{} != u) {
		utils.Log(2, "RegisterIdent", "requested user existing: "+u.Name)
		return errors.New("user already registred")
	} else {
		utils.Log(2, "RegisterIdent", "error checking user: ")
		return errors.New("error checking user:" + err.Error())
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
