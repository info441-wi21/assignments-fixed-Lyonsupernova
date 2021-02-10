package handlers

import (
	"assignments-fixed-Lyonsupernova/servers/gateway/models/users"
	"assignments-fixed-Lyonsupernova/servers/gateway/sessions"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

// UsersHandler handles the request
func (ch *ContextHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "The request method is not legal", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		log.Printf("The request body must be in JSON")
		http.Error(w, "The request body must be in JSON", http.StatusUnsupportedMediaType)
		return
	}
	// usr type users.NewUser
	newUsr := &users.NewUser{}
	err := json.NewDecoder(r.Body).Decode(newUsr)
	if err != nil {
		log.Printf("error decoding JSON: %v\n", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// any of the validation rules fail
	if err := newUsr.Validate(); err != nil {
		log.Printf("validate user fails")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// set usr as user from new user
	usr, err := newUsr.ToUser()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	insertUsr, err := ch.UserStore.Insert(usr)
	if err != nil {
		log.Printf("User database insert error")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// begin a session
	sessionState := SessionState{BeginDate: time.Now(), User: insertUsr}
	_, err = sessions.BeginSession(ch.SessionID, ch.SessionStore, sessionState, w)
	if err != nil {
		log.Printf("Begin session error")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// writing response to the clients
	w.Header().Add("Content-Type", "application/json")
	// status code 201
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(insertUsr); err != nil {
		log.Printf("User profile cannot be encoded in JSON format")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	json, _ := json.Marshal(insertUsr)
	w.Write([]byte(json))
}

// SpecificUserHandler authenticate the user and get the seesion state
func (ch *ContextHandler) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// authenticate the user and get the session state
	_, err := sessions.GetSessionID(r, ch.SessionID)
	if err != nil {
		http.Error(w, "current sessionID is not valid", http.StatusUnauthorized)
		return
	}

	sessionState := &SessionState{}
	_, err = sessions.GetState(r, ch.SessionID, ch.SessionStore, sessionState)

	if err != nil {
		http.Error(w, "current session state is not valid", http.StatusUnauthorized)
		return
	}

	base := filepath.Base(r.URL.Path)
	if r.Method == http.MethodGet {
		// get the user profile from the url path last element
		var usr *users.User
		var err error
		if base == "me" {
			usr, err = ch.UserStore.GetByID(sessionState.User.ID)
		} else {
			userID, err2 := strconv.ParseInt(base, 10, 64)
			if err2 != nil {
				http.Error(w, "User ID format not valid", http.StatusNotFound)
			}
			usr, err = ch.UserStore.GetByID(userID)
		}
		if err != nil {
			http.Error(w, "User ID not found in session store", http.StatusNotFound)
			return
		}
		// writing response to the clients
		w.Header().Set("Content-Type", "application/json")
		// status code 200
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(usr); err != nil {
			log.Printf("User profile cannot be encoded in JSON format")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json, _ := json.Marshal(usr)
		w.Write([]byte(json))
	} else if r.Method == http.MethodPatch {
		if base != "me" {
			userID, err := strconv.ParseInt(base, 10, 64)
			if err != nil {
				http.Error(w, "User ID format not valid", http.StatusBadRequest)
				return
			}
			_, err = ch.UserStore.GetByID(userID)
			if err != nil {
				log.Printf("The user is not authenticated")
				http.Error(w, "User ID not found in session store", http.StatusForbidden)
				return
			}
		}
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			http.Error(w, "The request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// update the current user profile with new profile
		newUsr := &users.Updates{}
		err := json.NewDecoder(r.Body).Decode(newUsr)
		if err != nil {
			log.Printf("error decoding JSON: %v\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		ch.UserStore.Update(sessionState.User.ID, newUsr)
		// writing response to the clients
		w.Header().Set("Content-Type", "application/json")
		// status code 200
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(newUsr); err != nil {
			log.Printf("User profile cannot be encoded in JSON format")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json, _ := json.Marshal(newUsr)
		w.Write([]byte(json))
	} else {
		http.Error(w, "Error request method", http.StatusMethodNotAllowed)
	}
}

// SessionsHandler handles ...
func (ch *ContextHandler) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			log.Printf("The request body must be in JSON")
			http.Error(w, "Bad request", http.StatusUnsupportedMediaType)
			return
		}
		userCredential := &users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(userCredential); err != nil {
			log.Printf("The request body cannot be encoded")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// find the user profile
		user, err := ch.UserStore.GetByEmail(userCredential.Email)
		if err != nil {
			log.Printf("The user's profile cannot be found")
			http.Error(w, "Profiel not found", http.StatusUnauthorized)
			return
		}

		// authentiicate
		if err := user.Authenticate(userCredential.Password); err != nil {
			log.Printf("The user's credential is not authentiated")
			http.Error(w, "Credential not authorized", http.StatusUnauthorized)
			return
		}

		// begin a new sessison
		sesssionState := SessionState{BeginDate: time.Now(), User: user}
		_, err = sessions.BeginSession(ch.SessionID, ch.SessionStore, sesssionState, w)
		if err != nil {
			log.Printf("Session cannot begin")
			http.Error(w, "Session cannot begin", http.StatusBadRequest)
			return
		}
		// writing response to the clients
		w.Header().Set("Content-Type", "application/json")
		// status code 201
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("User profile cannot be encoded in JSON format")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		json, _ := json.Marshal(user)
		w.Write([]byte(json))
	} else {
		http.Error(w, "Error status method: only accept Post", http.StatusMethodNotAllowed)
	}
}

// SpecificSessionHandler ...
func (ch *ContextHandler) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		base := filepath.Base(r.URL.Path)
		if base != "mine" {
			log.Printf("the last path segment does not equal \"mine\"")
			http.Error(w, "Session error", http.StatusForbidden)
			return
		}

		// end session
		sessions.EndSession(r, ch.SessionID, ch.SessionStore)

		// print out sign out information
		w.Write([]byte("signed out"))
	} else {
		http.Error(w, "Error status method: only accept Post", http.StatusMethodNotAllowed)
	}
}
