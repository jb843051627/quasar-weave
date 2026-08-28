package handler

import (
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
)

func (rt *Router) handleNotes(w http.ResponseWriter, r *http.Request, suffix string) {
	if len(splitPath(suffix)) != 0 {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodGet {
		items, err := rt.lab.ListNotes(r.Context(), r.URL.Query().Get("observation_id"), r.URL.Query().Get("alert_id"))
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
	var input model.NoteInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	item, err := rt.lab.AddNote(r.Context(), input)
	if err != nil {
		writeError(w, service.StatusCode(err), err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}
