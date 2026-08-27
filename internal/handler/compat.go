package handler

import "net/http"

func (h *Router) healthz(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Router) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Router) handleReadyz(w http.ResponseWriter, r *http.Request) {
	h.ready(w, r)
}

func (h *Router) handleHealth(w http.ResponseWriter, r *http.Request) {
	h.health(w, r)
}

func (h *Router) handleAntennas(w http.ResponseWriter, r *http.Request, suffix string) {
	h.antennas(w, r, suffix)
}
