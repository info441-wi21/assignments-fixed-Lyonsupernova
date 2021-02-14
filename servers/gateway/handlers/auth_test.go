package handlers

import (
	"assignments-fixed-Lyonsupernova/servers/gateway/models/users"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test UsersHandler handler
func TestUsersHandler(t *testing.T) {

	// Check different methods
	invalidMethods := [8]string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS",
		"TRACE", "PATCH"}

	for _, method := range invalidMethods {
		req, _ := http.NewRequest(method, "/v1/users", nil)
		rr := httptest.NewRecorder()
		UsersHandler(rr, req)
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("UserHandler accpet wrong methods %s", method)
		}
	}

	cases := []struct {
		returnStatus int
		contenType   string
		userFile     *users.NewUser
	}{
		{
			http.StatusUnsupportedMediaType,
			"application/x-www-form-urlencoded",
			nil,
		},
		{
			http.StatusBadRequest,
			"application/json",
			nil,
		},
		{
			http.StatusBadRequest,
			"application/json",
			&users.NewUser{
				Email: "aabbcc",
			},
		},
		{
			http.StatusCreated,
			"application/json",
			&users.NewUser{
				Email:        "aabbcc",
				Password:     "abc",
				PasswordConf: "notsurewhatisthis",
				UserName:     "userOne",
				FirstName:    "U",
				LastName:     "Ser",
			},
		},
	}

	for _, c := range cases {
		reqBody, _ := json.Marshal(c.userFile)
		req, _ := http.NewRequest("POST", "/v1/users", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", c.contenType)
		rr := httptest.NewRecorder()
		UsersHandler(rr, req)
		if status := rr.Code; status != c.returnStatus {
			t.Errorf("Instead of status %s, UserHandler response with %s http status",
				c.returnStatus, status)
		}
		var user *users.User
		json.Unmarshal(rr.Body.Bytes(), &user)
		if user.PassHash != nil || user.Email != "" {
			t.Errorf("Response with unnecessary PassHash or Email")
		}
	}
}

// Test SpecificUserHandler
func TestSpecificUserHandler(t *testing.T) {

	// Create a dummy user to test for
	user := &users.NewUser{
		Email:        "secondEmail@gmail.com",
		Password:     "123",
		PasswordConf: "456",
		UserName:     "usertwo",
		FirstName:    "bat",
		LastName:     "man",
	}

	reqBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/v1/users", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Errorf("Did register the user successfully")
	}
	rr := httptest.NewRecorder()
	UsersHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		fmt.Errorf("Did register the user successfully")
	}
	var response *users.User
	json.Unmarshal(rr.Body.Bytes(), &user)
	userID := response.ID

	//TODO: need to check if it is authenticate

	// start testing cases
	cases := []struct {
		id               string
		expectedResponse int
	}{
		{
			"iaminvalidID",
			http.StatusNotFound,
		},
		{
			string(userID),
			http.StatusOK,
		},
		{
			"me",
			http.StatusOK,
		},
	}
}
