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
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
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

	db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Erorr connecting to db %v", err)
	}
	defer db.Close()

	_, err = sessions.NewSessionID("somerandomkey")
	if err != nil {
		fmt.Printf("Error generating sessionID")
	}
	userStore := &users.DummyMySQLStore{}
	//ctxSess := sessions.NewMemStore(100, 10)
	contextHandler := &ContextHandler{
		SessionID:    "ramdom",
		SessionStore: sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), time.Hour),
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
		validateSessID   bool
		update           *users.Updates
		user             *users.User
	}{
		{
			"SpecificUserHandler1",
			"notme",
			"GET",
			"",
			http.StatusNotFound,
			true,
			&users.Updates{},
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler2",
			"51231516",
			"GET",
			"",
			http.StatusOK,
			true,
			&users.Updates{},
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler3",
			"me",
			"GET",
			"",
			http.StatusOK,
			true,
			&users.Updates{},
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler4",
			"1231709",
			"PATCH",
			"application/json",
			http.StatusForbidden,
			false,
			&users.Updates{},
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler5",
			"me",
			"PATCH",
			"application/notjson",
			http.StatusUnsupportedMediaType,
			true,
			&users.Updates{},
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler6",
			"me",
			"PATCH",
			"application/json",
			http.StatusOK,
			true,
			&users.Updates{
				FirstName: "new",
				LastName:  "man",
			},
			&users.User{
				ID: 123456,
			},
		},
		{
			"SpecificUserHandler7",
			"me",
			"POST",
			"application/json",
			http.StatusMethodNotAllowed,
			true,
			&users.Updates{},
			&users.User{},
		},
	}

	for _, c := range cases {
		log.Printf("Testing %s ...", c.sampleID)
		body := []byte("{json:json}")
		if c.method == "PATCH" {
			body, err = json.Marshal(c.update)
			if err != nil {
				t.Errorf("Error json encoding")
			}
		}
		handler := http.HandlerFunc(contextHandler.SpecificUserHandler)
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(c.method, "/v1/users/"+c.id, bytes.NewReader(body))
		sess := &SessionState{
			BeginDate: time.Now(),
			User:      c.user,
		}
		sid, err := sessions.BeginSession(contextHandler.SessionID, contextHandler.SessionStore, sess, rr)
		if err != nil {
			t.Errorf("Error beginning sessions, %v", err)
		}
		if c.validateSessID == false {
			req.Header.Set("Authorization", "invalidID")
		} else {
			req.Header.Set("Authorization", sid.String())
		}
		handler.ServeHTTP(rr, req)
		err = contextHandler.SessionStore.Save(sid, sess)
		if err != nil {
			t.Errorf("session save went wrong")
		}
		// checks if it returns with a correct status code
		if status := rr.Code; status != c.expectedResponse {
			t.Errorf("Instead of status %d, handler response with %d http status, test %s",
				c.expectedResponse, status, c.sampleID)
		}
		log.Printf("%s Passed", c.sampleID)
	}
}

func TestSessionsHandler(t *testing.T) {
	db, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("Erorr connecting to db %v", err)
	}
	defer db.Close()
	_, err = sessions.NewSessionID("somerandomkey")
	if err != nil {
		fmt.Printf("Error generating sessionID")
	}
	userStore := &users.DummyMySQLStore{}
	//ctxSess := sessions.NewMemStore(100, 10)
	contextHandler := &ContextHandler{
		SessionID:    "ramdom",
		SessionStore: sessions.NewRedisStore(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}), time.Hour),
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
		user             *users.User
		hasCred          int
	}{
		{
			"SessionsHandler1",
			"GET",
			"",
			http.StatusMethodNotAllowed,
			&users.Credentials{},
			&users.User{},
			0,
		},
		{
			"SessionsHandler2",
			"POST",
			"application/x-www-form-urlencoded",
			http.StatusUnsupportedMediaType,
			&users.Credentials{},
			&users.User{},
			0,
		},
		{
			"SessionsHandler3",
			"POST",
			"application/json",
			http.StatusCreated,
			credential,
			&users.User{},
			1,
		},
	}

	for _, c := range cases {
		log.Printf("Testing %s ...", c.sampleID)

		handler := http.HandlerFunc(contextHandler.SessionsHandler)
		body := []byte("{json:json}")
		if c.credentails.Email != "" && c.credentails.Password != "" {
			body, _ = json.Marshal(c.credentails)
		}
		req, _ := http.NewRequest(c.method, "/v1/sessions", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", c.contentType)
		if c.hasCred == 1 {
			c.user.PassHash, _ = bcrypt.GenerateFromPassword(c.user.PassHash, 13)
		}
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
