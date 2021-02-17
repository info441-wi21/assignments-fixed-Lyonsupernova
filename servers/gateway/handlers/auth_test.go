package handlers

import (
	"assignments-fixed-Lyonsupernova/servers/gateway/models/users"
	"assignments-fixed-Lyonsupernova/servers/gateway/sessions"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test UsersHandler handler
func TestUsersHandler(t *testing.T) {

	user := &users.DummyMySQLStore{}
	sess := sessions.NewMemStore(100, 10)
	contextHandler := &ContextHandler{
		SessionID:    "ramdom",
		SessionStore: sess,
		UserStore:    user,
	}
	// Check different methods
	invalidMethods := [8]string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS",
		"TRACE", "PATCH"}

	for _, method := range invalidMethods {
		req, _ := http.NewRequest(method, "/v1/users", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(contextHandler.UsersHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("UserHandler accpet wrong methods %s", method)
		}
	}
	log.Printf("Finish methods tests")

	cases := []struct {
		sampleID         string
		contenType       string
		expectedResponse int
		userFile         *users.NewUser
	}{
		{
			"UsersHandler1",
			"application/x-www-form-urlencoded",
			http.StatusUnsupportedMediaType,
			&users.NewUser{},
		},
		{
			"UsersHandler2",
			"application/json",
			http.StatusBadRequest,
			&users.NewUser{},
		},
		{
			"UsersHandler3",
			"application/json",
			http.StatusBadRequest,
			&users.NewUser{
				Email: "aabbcc",
			},
		},
		{
			"UsersHandler4",
			"application/json",
			http.StatusCreated,
			&users.NewUser{
				Email:        "aabbcc@gmail.com",
				Password:     "abcefg",
				PasswordConf: "abcefg",
				UserName:     "userOne",
				FirstName:    "U",
				LastName:     "Ser",
			},
		},
	}

	for _, c := range cases {
		log.Printf("Testing %s ...", c.sampleID)
		reqBody := new(bytes.Buffer)
		bufEncode := json.NewEncoder(reqBody)
		bufEncode.Encode(c.userFile)
		//reqBody, _ := json.Marshal(c.userFile)
		req, _ := http.NewRequest("POST", "/v1/users", bytes.NewReader(reqBody.Bytes()))
		req.Header.Set("Content-Type", c.contenType)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(contextHandler.UsersHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, handler response with %d http status",
				c.expectedResponse, status)
		}
		var user *users.User
		user = new(users.User)
		json.Unmarshal(rr.Body.Bytes(), &user)
		if user.PassHash != nil || user.Email != "" {
			t.Errorf("Response with unnecessary PassHash or Email")
		}
		log.Printf("%s Passed", c.sampleID)
	}
}

// Test SpecificUserHandler
func TestSpecificUserHandler(t *testing.T) {

	_, err := sessions.NewSessionID("random")
	if err != nil {
		log.Printf("Didn't generate sessionsID")
	}
	userStore := &users.DummyMySQLStore{}
	sess := sessions.NewMemStore(1000, 100)
	contextHandler := &ContextHandler{
		SessionID:    "ramdom",
		SessionStore: sess,
		UserStore:    userStore,
	}

	/*
		// Create a dummy user to test for
		user := &users.NewUser{
			Email:        "secondEmail@gmail.com",
			Password:     "123535",
			PasswordConf: "123535",
			UserName:     "usertwo",
			FirstName:    "bat",
			LastName:     "man",
		}
		rr := httptest.NewRecorder()
		registerUser(user, contextHandler, rr)
		// Retrieve user id from the response for further testings
		var response *users.User
		response = new(users.User)
		json.Unmarshal(rr.Body.Bytes(), &user)
		userID := response.ID
	*/

	// start testing cases
	cases := []struct {
		sampleID         string
		id               string
		method           string
		contentType      string
		expectedResponse int
		user             *users.User
	}{
		{
			"SpecificUserHandler1",
			"iaminvalidID",
			http.MethodGet,
			"",
			http.StatusUnauthorized,
			&users.User{
				ID: 492012,
			},
		},
		{
			"SpecificUserHandler2",
			"123456",
			http.MethodGet,
			"",
			http.StatusOK,
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler3",
			"23456",
			http.MethodGet,
			"",
			http.StatusOK,
			&users.User{
				ID: 23456,
			},
		},
		{
			"SpecificUserHandler4",
			"notme",
			http.MethodPatch,
			"application/json",
			http.StatusForbidden,
			&users.User{
				ID: 23456,
			},
		},
		{
			"SpecificUserHandler5",
			"989555",
			http.MethodPatch,
			"application/notjson",
			http.StatusUnsupportedMediaType,
			&users.User{
				ID: 989555,
			},
		},
		{
			"SpecificUserHandler6",
			"896967",
			http.MethodPatch,
			"application/json",
			http.StatusOK,
			&users.User{
				ID: 896967,
			},
		},
		{
			"SpecificUserHandler7",
			"me",
			http.MethodPost,
			"application/json",
			http.StatusMethodNotAllowed,
			&users.User{},
		},
	}

	for _, c := range cases {
		log.Printf("Testing %s ...", c.sampleID)
		req, err := http.NewRequest(c.method, "/v1/users/"+c.id, nil)
		if err != nil {
			log.Printf("Error generating Newrequest")
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(contextHandler.SpecificUserHandler)
		handler.ServeHTTP(rr, req)
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, handler response with %d http status",
				c.expectedResponse, status)
		}
		log.Printf("%s Passed", c.sampleID)
	}
}

func TestSessionsHandler(t *testing.T) {
	userStore := &users.DummyMySQLStore{}
	sess := sessions.NewMemStore(100, 10)
	contextHandler := &ContextHandler{
		SessionID:    "ramdom",
		SessionStore: sess,
		UserStore:    userStore,
	}

	// Create a dummy user to test for
	user := &users.NewUser{
		Email:        "thridEmail@gmail.com",
		Password:     "513451",
		PasswordConf: "513451",
		UserName:     "userThree",
		FirstName:    "super",
		LastName:     "man",
	}

	credential := &users.Credentials{
		Email:    "thridEmail@gmail.com",
		Password: "513451",
	}

	rr := httptest.NewRecorder()
	registerUser(user, contextHandler, rr)

	// Start test cases
	cases := []struct {
		sampleID         string
		method           string
		contentType      string
		expectedResponse int
		credentails      *users.Credentials
	}{
		{
			"SessionsHandler1",
			"GET",
			"",
			http.StatusMethodNotAllowed,
			nil,
		},
		{
			"SessionsHandler2",
			"POST",
			"application/x-www-form-urlencoded",
			http.StatusUnsupportedMediaType,
			nil,
		},
		{
			"SessionsHandler3",
			"POST",
			"application/json",
			http.StatusCreated,
			credential,
		},
	}

	for _, c := range cases {
		log.Printf("Testing %s ...", c.sampleID)
		reqBody := new(bytes.Buffer)
		bufEncode := json.NewEncoder(reqBody)
		bufEncode.Encode(c.credentails)
		req, _ := http.NewRequest(c.method, "/v1/sessions", bytes.NewReader(reqBody.Bytes()))
		req.Header.Set("contentType", c.contentType)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(contextHandler.SessionsHandler)
		handler.ServeHTTP(rr, req)
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, handler response with %d http status",
				c.expectedResponse, status)
		}
		log.Printf("%s Passed", c.sampleID)
	}
}

func TestSpecificSessionHandler(t *testing.T) {
	userStore := &users.DummyMySQLStore{}
	sess := sessions.NewMemStore(100, 10)
	contextHandler := &ContextHandler{
		SessionID:    "ramdom",
		SessionStore: sess,
		UserStore:    userStore,
	}

	// Create a dummy user to test for
	user := &users.NewUser{
		Email:        "forthEmail@gmail.com",
		Password:     "1236135",
		PasswordConf: "1236135",
		UserName:     "userFour",
		FirstName:    "spider",
		LastName:     "man",
	}

	rr := httptest.NewRecorder()
	registerUser(user, contextHandler, rr)

	cases := []struct {
		sampleID         string
		method           string
		lastURL          string
		expectedResponse int
	}{
		{
			"SpecificSessionHandler1",
			"GET",
			"mine",
			http.StatusMethodNotAllowed,
		},
		{
			"SpecificSessionHandler2",
			"DELETE",
			"me",
			http.StatusForbidden,
		},
		{
			"SpecificSessionHandler3",
			"DELET",
			"mine",
			http.StatusOK,
		},
	}

	for _, c := range cases {
		log.Printf("Testing %s ...", c.sampleID)
		req, _ := http.NewRequest(c.method, "/v1/sessions/"+c.lastURL, nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(contextHandler.SpecificSessionHandler)
		handler.ServeHTTP(rr, req)
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, handler response with %d http status",
				c.expectedResponse, status)
		}
		log.Printf("%s Passed", c.sampleID)
	}
}

// A helper function to register for a dummy user
func registerUser(user *users.NewUser, contextHandler *ContextHandler,
	rr *httptest.ResponseRecorder) {
	reqBody := new(bytes.Buffer)
	bufEncode := json.NewEncoder(reqBody)
	bufEncode.Encode(user)
	//reqBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewReader(reqBody.Bytes()))
	handler := http.HandlerFunc(contextHandler.UsersHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		fmt.Errorf("Didn't register the user successfully")
	}
}

// A helper function to generate a user
func dumyNewUser() *users.NewUser {
	return &users.NewUser{
		Email:        "ForthEmail@gmail.com",
		Password:     "123535",
		PasswordConf: "123535",
		UserName:     "usertwo",
		FirstName:    "bat",
		LastName:     "man",
	}
}

func dumySpecifiedUser(id int64) *users.User {
	return &users.User{
		ID:        id,
		Email:     "Fifthemail@gmail.com",
		PassHash:  []byte("password123"),
		UserName:  "Dragon",
		FirstName: "Dragon",
		LastName:  "Man",
		PhotoURL:  "SomeURL",
	}
}