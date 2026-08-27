package handler

import (
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
)

func (rt *Router) handleGates(w http.ResponseWriter, r *http.Request, suffix string) {
	parts := splitPath(suffix)
	if len(parts) == 0 {
		if r.Method == http.MethodGet {
			active := r.URL.Query().Get("active") == "true"
			items, err := rt.lab.ListGates(r.Context(), active)
			if err != nil {
				writeError(w, service.StatusCode(err), err)
				return
			}
			writeOK(w, items)
			return
		}
		if !method(w, r, http.MethodPost) {
			return
		}
		var input model.GateInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := rt.lab.CreateGate(r.Context(), input)
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	if len(parts) == 2 && parts[1] == "active" {
		if !method(w, r, http.MethodPost) {
			return
		}
		var body struct {
			Active bool `json:"active"`
		}
		if err := decodeJSON(r, &body); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := rt.lab.SetGateActive(r.Context(), parts[0], body.Active)
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	if len(parts) == 1 {
		if !method(w, r, http.MethodGet) {
			return
		}
		item, err := rt.lab.GetGate(r.Context(), parts[0])
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	http.NotFound(w, r)
}
