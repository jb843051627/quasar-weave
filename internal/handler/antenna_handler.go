package handler

import (
	"net/http"
	"strings"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (h *Router) antennas(w http.ResponseWriter, r *http.Request, suffix string) {
	suffix = strings.Trim(suffix, "/")
	parts := splitPath(suffix)
	if len(parts) == 0 {
		if r.Method == http.MethodGet {
			items, err := h.lab.ListAntennas(r.Context())
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, items)
			return
		}
		if !method(w, r, http.MethodPost) {
			return
		}
		var input model.AntennaInput
		if err := decodeBody(r, &input); err != nil {
			writeError(w, err)
			return
		}
		item, err := h.lab.RegisterAntenna(r.Context(), input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	id := parts[0]
	if len(parts) == 2 && parts[1] == "heartbeat" {
		if !method(w, r, http.MethodPost) {
			return
		}
		item, err := h.lab.Heartbeat(r.Context(), id, model.AntennaReady)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	if len(parts) == 2 && parts[1] == "enabled" {
		if !method(w, r, http.MethodPatch) {
			return
		}
		var body struct {
			Enabled bool `json:"enabled"`
		}
		if err := decodeBody(r, &body); err != nil {
			writeError(w, err)
			return
		}
		item, err := h.lab.SetAntennaEnabled(r.Context(), id, body.Enabled)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	if len(parts) == 1 && r.Method == http.MethodGet {
		item, err := h.lab.GetAntenna(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
		return
	}
	http.NotFound(w, r)
}
