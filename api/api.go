package api

import (
	"encoding/json"
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

	type Test struct {
		MsgType string `json:"MsgType"`
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.Log(3, "TestApi() ", "reading reqBody : "+err.Error())
	}

	m := Test{}
	errunmarsh := json.Unmarshal(reqBody, &m)
	if errunmarsh != nil {
		utils.Log(3, "TestApi() ", "unmarshal reqBody : ")

	}

	switch m.MsgType {
	case "test":
		utils.Log(1, "TestApi() ", "got messagetype Test")
	case "RegisterRequest":
		utils.Log(1, "TestApi() ", "got messagetype RegisterRequest")
	default:
		utils.Log(3, "TestApi() ", "got unknown messagetype: "+m.MsgType)
	}

	w.Write([]byte("{TEST-OK-Answer}"))

}

// ToDo: find correct user and do the response
func RegisterIdent(w http.ResponseWriter, r *http.Request) {
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
