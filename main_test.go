package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wyrdnixx/KBZOZeit/api"
	"github.com/wyrdnixx/KBZOZeit/database"
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

func TestUser(t *testing.T) {
	expected := models.User{
		Name:    "EnabledExampleUser",
		Enabled: 0,
	}

	errAddUser := database.AddUser("EnabledExampleUser")
	if errAddUser != nil {
		t.Errorf("AddUser got error: got %v", errAddUser)
	}

	errDisableUser := database.DisableUser("EnabledExampleUser")
	if errDisableUser != nil {
		t.Errorf("DisableUser got error: got %v", errDisableUser)
	}

	u, errFindUser := database.FindUser("EnabledExampleUser")
	if errFindUser != nil {
		t.Errorf("TestFindUser returned unexpected error: got %v want %v", errFindUser, expected)
	}
	if u != expected {
		t.Errorf("TestFindUser returned unexpected result: got %v want %v", u, expected)
	}

	users, err := database.GetUsers()
	if err != nil {
		t.Errorf("GetUsers error: %v", err)
	}
	if len(users.User) == 0 {
		t.Errorf("GetUsers got no users: %v", users)
	}

}
