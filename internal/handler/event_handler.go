package handler

import "net/http"

func (rt *Router) handleEvents(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	items, err := rt.lab.Events(r.Context(), r.URL.Query().Get("subject"), 100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeOK(w, items)
}
