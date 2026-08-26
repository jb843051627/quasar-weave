package handler

import "net/http"

func (h *Router) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		method(w, r, http.MethodGet)
		return
	}
	item, err := h.lab.Health(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *Router) ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		method(w, r, http.MethodGet)
		return
	}
	if err := h.lab.Readiness(r.Context()); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (h *Router) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		method(w, r, http.MethodGet)
		return
	}
	items, err := h.lab.Events(r.Context(), r.URL.Query().Get("subject"), 100)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
