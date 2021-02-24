package main

import (
	"assignments-fixed-Lyonsupernova/servers/summary/handlers"
	"log"
	"net/http"
	"os"
)

// main for the summary api
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary/", handlers.SummaryHandler)
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux)) // http for inner network will be suffcient
}
