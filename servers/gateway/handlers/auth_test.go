package handlers

import (
	"assignments-fixed-Lyonsupernova/servers/gateway/models/users"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Test UsersHandler handler
func TestUsersHandler(t *testing.T) {

	contextHandler := &ContextHandler{}
	// Check different methods
	invalidMethods := [8]string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS",
		"TRACE", "PATCH"}

	for _, method := range invalidMethods {
		req, _ := http.NewRequest(method, "/v1/users", nil)
		rr := httptest.NewRecorder()
		contextHandler.UsersHandler(rr, req)
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("UserHandler accpet wrong methods %s", method)
		}
	}

	cases := []struct {
		expectedResponse int
		contenType       string
		userFile         *users.NewUser
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
		contextHandler.UsersHandler(rr, req)
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, UserHandler response with %d http status",
				c.expectedResponse, status)
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

	contextHandler := &ContextHandler{}

	// Create a dummy user to test for
	user := &users.NewUser{
		Email:        "secondEmail@gmail.com",
		Password:     "123",
		PasswordConf: "456",
		UserName:     "usertwo",
		FirstName:    "bat",
		LastName:     "man",
	}

	rr := httptest.NewRecorder()
	registerUser(user, contextHandler, rr)

	// Retrieve user id from the response for further testings
	var response *users.User
	json.Unmarshal(rr.Body.Bytes(), &user)
	userID := response.ID

	// start testing cases
	cases := []struct {
		id               string
		method           string
		contentType      string
		expectedResponse int
	}{
		{
			"iaminvalidID",
			"GET",
			"",
			http.StatusNotFound,
		},
		{
			fmt.Sprint(userID),
			"GET",
			"",
			http.StatusOK,
		},
		{
			"me",
			"GET",
			"",
			http.StatusOK,
		},
		{
			"notme",
			"PATCH",
			"application/json",
			http.StatusForbidden,
		},
		{
			"me",
			"PATCH",
			"application/notjson",
			http.StatusUnsupportedMediaType,
		},
		{
			"me",
			"PATCH",
			"application/json",
			http.StatusOK,
		},
		{
			"me",
			"POST",
			"application/json",
			http.StatusMethodNotAllowed,
		},
	}

	for _, c := range cases {
		req, _ := http.NewRequest(c.method, "/v1/users/"+c.id, nil)
		rr := httptest.NewRecorder()
		contextHandler.UsersHandler(rr, req)
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, UserHandler response with %d http status",
				c.expectedResponse, status)
		}
		// checks if it returns a complete json
		var values reflect.Value
		if c.method == "GET" {
			var user *users.User
			json.Unmarshal(rr.Body.Bytes(), &user)
			values = reflect.ValueOf(user)
		} else {
			var update *users.Updates
			json.Unmarshal(rr.Body.Bytes(), &update)
			values = reflect.ValueOf(update)
		}
		for i := 0; i < values.NumField(); i++ {
			if values.Field(i).Interface() == nil {
				t.Errorf("Response json does not have %s", fmt.Sprint(values.Field(i)))
			}
		}
	}
}

func TestSessionsHandler(t *testing.T) {
	contextHandler := &ContextHandler{}

	// Create a dummy user to test for
	user := &users.NewUser{
		Email:        "thridEmail@gmail.com",
		Password:     "123",
		PasswordConf: "456",
		UserName:     "userThree",
		FirstName:    "super",
		LastName:     "man",
	}

	credential := &users.Credentials{
		Email:    "thridEmail@gmail.com",
		Password: "123",
	}

	rr := httptest.NewRecorder()
	registerUser(user, contextHandler, rr)

	// Start test cases
	cases := []struct {
		method           string
		contentType      string
		expectedResponse int
		credentails      *users.Credentials
	}{
		{
			"GET",
			"",
			http.StatusMethodNotAllowed,
			nil,
		},
		{
			"POST",
			"application/x-www-form-urlencoded",
			http.StatusUnsupportedMediaType,
			nil,
		},
		{
			"POST",
			"application/json",
			http.StatusCreated,
			credential,
		},
	}

	for _, c := range cases {
		var req *http.Request
		if c.credentails != nil {
			reqBody, _ := json.Marshal(c.credentails)
			req, _ = http.NewRequest(c.method, "/v1/sessions", bytes.NewReader(reqBody))
		} else {
			req, _ = http.NewRequest(c.method, "/v1/sessions", nil)
		}
		rr := httptest.NewRecorder()
		contextHandler.UsersHandler(rr, req)
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, UserHandler response with %d http status",
				c.expectedResponse, status)
		}
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	contextHandler := &ContextHandler{}

	// Create a dummy user to test for
	user := &users.NewUser{
		Email:        "forthEmail@gmail.com",
		Password:     "123",
		PasswordConf: "456",
		UserName:     "userFour",
		FirstName:    "spider",
		LastName:     "man",
	}

	rr := httptest.NewRecorder()
	registerUser(user, contextHandler, rr)

	cases := []struct {
		method           string
		lastURL          string
		expectedResponse int
	}{
		{
			"GET",
			"mine",
			http.StatusMethodNotAllowed,
		},
		{
			"DELETE",
			"me",
			http.StatusForbidden,
		},
		{
			"DELET",
			"mine",
			http.StatusOK,
		},
	}

	for _, c := range cases {
		req, _ := http.NewRequest(c.method, "/v1/sessions/"+c.lastURL, nil)
		rr := httptest.NewRecorder()
		contextHandler.UsersHandler(rr, req)
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, UserHandler response with %d http status",
				c.expectedResponse, status)
		}
	}

}

// A helper function to register for a dummy user
func registerUser(user *users.NewUser, contextHandler *ContextHandler,
	rr *httptest.ResponseRecorder) {
	reqBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewReader(reqBody))
	contextHandler.UsersHandler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		fmt.Errorf("Didn't register the user successfully")
	}
}
