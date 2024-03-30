package main

import (
	"fmt"
	"net/http"

	"github.com/shaikzhafir/go-htmx-starter/internal/commons"
	"github.com/shaikzhafir/go-htmx-starter/internal/handlers"
	log "github.com/shaikzhafir/go-htmx-starter/internal/logging"
)

func main() {
	mux := http.NewServeMux()

	handlers := handlers.NewAPIHandler()
	// css and js files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// other paths
	mux.HandleFunc("/api", handlers.GetFakeData())
	// html files (also catch all paths)
	mux.Handle("/", http.FileServer(http.Dir("public")))

	log.Info("Server started on port %d", commons.DefaultPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", commons.DefaultPort), mux)
	if err != nil {
		log.Error("Error starting server %s", err.Error())
	}
}
