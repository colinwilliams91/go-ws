package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/colinwilliams91/go-ws.git/internal/handlers"
)

// TODO: main package pointing to GH is awkward? could just be "go-ws" for internal use?

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))

	return mux
}