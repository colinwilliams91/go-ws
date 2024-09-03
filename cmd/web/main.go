package main

import (
	"log"
	"net/http"

	"github.com/colinwilliams91/go-ws.git/internal/handlers"
)

// placeholder

func main() {
	var mux http.Handler = routes()

	log.Println("Starting channel listener...")

	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080...")

	_ = http.ListenAndServe(":8080", mux)
}