package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/shaikzhafir/go-htmx-starter/internal/commons"
	h "github.com/shaikzhafir/go-htmx-starter/internal/handlers"
	log "github.com/shaikzhafir/go-htmx-starter/internal/logging"
	store "github.com/shaikzhafir/go-htmx-starter/shared"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
)

// Config configures the main ServeMux.
type Config struct {
	ClientID     string
	ClientSecret string
}

func main() {
	config := initAuthStuff()
	mux := initServeMux(config)
	log.Info("Server started on port %d", commons.DefaultPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", commons.DefaultPort), mux)
	if err != nil {
		log.Error("Error starting server %s", err.Error())
	}
}

func initServeMux(config *Config) *http.ServeMux {
	mux := http.NewServeMux()
	store := store.NewSessionStore()
	handlers := h.NewAPIHandler()
	html := h.NewHTMLHandler(store)
	auth := h.NewAuthHandler(store)
	// css and js files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// other paths
	mux.HandleFunc("/api", handlers.GetFakeData())
	mux.HandleFunc("/logout", auth.Logout())
	// 1. Register Login and Callback handlers
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "http://localhost:8080/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	mux.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, auth.LoginCallback(), nil)))
	mux.Handle("/profile", html.Profile())
	// html files (also catch all paths)
	mux.Handle("/", html.Index())
	return mux
}

func initAuthStuff() *Config {
	config := &Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}
	// allow consumer credential flags to override config fields
	clientID := flag.String("client-id", "", "Google Client ID")
	clientSecret := flag.String("client-secret", "", "Google Client Secret")
	flag.Parse()
	if *clientID != "" {
		config.ClientID = *clientID
	}
	if *clientSecret != "" {
		config.ClientSecret = *clientSecret
	}
	if config.ClientID == "" {
		log.Fatal("Missing Google Client ID")
	}
	if config.ClientSecret == "" {
		log.Fatal("Missing Google Client Secret")
	}
	return config
}
