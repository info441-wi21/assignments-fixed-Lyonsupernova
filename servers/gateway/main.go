package main

import (
	"assignments-fixed-Lyonsupernova/servers/gateway/handlers"
	"assignments-fixed-Lyonsupernova/servers/gateway/models/users"
	"assignments-fixed-Lyonsupernova/servers/gateway/sessions"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

//main is the main entry point for the server
func main() {

	sessionID := os.Getenv("SESSIONKEY")
	redisAddr := os.Getenv("REDISADDR")
	TLSKEY := os.Getenv("TLSKEY")
	TLSCERT := os.Getenv("TLSCERT")
	messageAddr := strings.Split(os.Getenv("MESSAGESADDR"), ",")
	summaryAddr := strings.Split(os.Getenv("SUMMARYADDR"), ",")
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	DSN := os.Getenv("DSN")
	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}
	redisDB := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	sessionStore := sessions.NewRedisStore(redisDB, time.Hour)

	userStore, err := users.NewMySQLStore(DSN)
	if err != nil {
		log.Printf("Unable to open database mysql %v", err)
	}

	contextHandler := &handlers.ContextHandler{
		SessionID:    sessionID,
		SessionStore: sessionStore,
		UserStore:    userStore,
	}

	messageDirector := func(r *http.Request) {
		if len(r.Header.Get("X-User")) == 0 {
			r.Header.Add("X-user", "")
		}
		auth := r.Header.Get("Authorization")
		if len(auth) != 0 {
			sessID := sessions.SessionID(strings.TrimPrefix(auth, "Bearer "))
			sessState := &handlers.SessionState{}
			err := contextHandler.SessionStore.Get(sessID, sessState)
			if err == nil {
				r.Header.Set("X-User", string(sessState.User.ID))
			} else {
				r.Header.Del("X-User")
			}
		}

		rand.Seed(time.Now().UnixNano())
		serveNum := rand.Intn(len(messageAddr))
		r.Host = messageAddr[serveNum]
		r.URL.Host = messageAddr[serveNum]
		r.URL.Scheme = "http"
	}
	messageProxy := &httputil.ReverseProxy{Director: messageDirector}

	summaryDirector := func(r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		serveNum := rand.Intn(len(summaryAddr))
		r.Host = summaryAddr[serveNum]
		r.URL.Host = summaryAddr[serveNum]
		r.URL.Scheme = "http"
	}
	summaryProxy := &httputil.ReverseProxy{Director: summaryDirector}

	mux := http.NewServeMux()
	log.Printf("server is listening at %s...", addr)
	mux.HandleFunc("/v1/users", contextHandler.UsersHandler)
	mux.HandleFunc("/v1/users/", contextHandler.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", contextHandler.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", contextHandler.SpecificSessionHandler)
	mux.Handle("/v1/channels", messageProxy)
	mux.Handle("/v1/channels/", messageProxy)
	mux.Handle("/v1/messages/", messageProxy)
	mux.Handle("/v1/summary/", summaryProxy)
	wrappedMux := handlers.NewHeaderHandler(mux)
	log.Fatal(http.ListenAndServeTLS(addr, TLSCERT, TLSKEY, wrappedMux))
}
