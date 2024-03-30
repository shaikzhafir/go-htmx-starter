package handlers

import "net/http"

func NewAPIHandler() *APIHandler {
	return &APIHandler{}
}

type APIHandler struct {
}

func (a *APIHandler) GetFakeData() http.HandlerFunc {
	// handle the request
	return func(w http.ResponseWriter, r *http.Request) {
		// get params from request
		// do some processing
		// return response
		name := r.FormValue("name")
		w.Write([]byte("wow this is fake data!, u wrote " + name + "in the request!"))
	}
}
