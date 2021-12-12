package api

import (
	"encoding/json"
	"log"
	"net/http"
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
