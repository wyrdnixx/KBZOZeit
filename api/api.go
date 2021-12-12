package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wyrdnixx/KBZOZeit/database"
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
		Msg   string `json:"msg"`
		Value bool   `json:"value"`
	}
	var SendTest Test
	SendTest.Msg = "Test"
	SendTest.Value = true

	response, err := json.Marshal(SendTest)
	if err != nil {
		log.Fatal("JSON Error")
	} else {
		w.Write(response)
	}

}

func AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	u, err := database.GetUsers()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		response, err := json.Marshal(u)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(response)
		}
	}

}
