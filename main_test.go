package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wyrdnixx/KBZOZeit/frontend/api"
	"github.com/wyrdnixx/KBZOZeit/frontend/database"
	"github.com/wyrdnixx/KBZOZeit/models"
)

func TestApi(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/TestApi", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.TestApi)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"msg":"Test","value":true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestFindUser(t *testing.T) {
	expected := models.User{
		Id:      1,
		Name:    "EnabledExampleUser",
		Enabled: 1,
	}
	u, err := database.FindUser("EnabledExampleUser")
	if err != nil {
		t.Errorf("TestFindUser returned unexpected error: got %v want %v", err, expected)
	}
	if u != expected {
		t.Errorf("TestFindUser returned unexpected result: got %v want %v", u, expected)
	}

}
