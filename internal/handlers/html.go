package handlers

import (
	"net/http"
	"text/template"

	"github.com/dghubble/sessions"
	"github.com/pkg/errors"
	log "github.com/shaikzhafir/go-htmx-starter/internal/logging"
)

type HTMLHandler interface {
	Index() http.HandlerFunc
	Profile() http.HandlerFunc
}

type htmlHandler struct {
	sessionStore sessions.Store[string]
}

type UserInfo struct {
	Name  string
	Email string
}

func NewHTMLHandler(s sessions.Store[string]) HTMLHandler {
	return &htmlHandler{
		sessionStore: s,
	}
}

func (h *htmlHandler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "./templates/home.page.html", nil)
	}
}

func (h *htmlHandler) Profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionStore.Get(r, "example-google-app")
		log.Info("Session %v", session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render(w, "./templates/profile.page.html", map[string]string{
			"Name":   session.Get("googleName"),
			"Email":  session.Get("googleEmail"),
			"Avatar": session.Get("googleAvatar"),
		})
	}
}

func render(w http.ResponseWriter, path string, data map[string]string) {
	tmpl, err := template.ParseFiles(path, "./templates/main.layout.html")
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to render html page").Error(), http.StatusInternalServerError)
		return
	}
	if data != nil {
		tmpl.ExecuteTemplate(w, "main", data)
		return
	}
	tmpl.Execute(w, nil)
}
