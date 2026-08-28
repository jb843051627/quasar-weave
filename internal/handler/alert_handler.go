package handler

import (
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
)

func (rt *Router) handleAlerts(w http.ResponseWriter, r *http.Request, suffix string) {
	parts := splitPath(suffix)
	if len(parts) == 0 {
		if !method(w, r, http.MethodGet) {
			return
		}
		filter := model.AlertFilter{ObservationID: r.URL.Query().Get("observation_id"), State: model.AlertState(r.URL.Query().Get("state")), Severity: r.URL.Query().Get("severity")}
		items, err := rt.lab.ListAlerts(r.Context(), filter)
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, items)
		return
	}
	if len(parts) == 1 {
		if !method(w, r, http.MethodGet) {
			return
		}
		item, err := rt.lab.GetAlert(r.Context(), parts[0])
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	if len(parts) == 2 && (parts[1] == "ack" || parts[1] == "resolve") {
		if !method(w, r, http.MethodPost) {
			return
		}
		var body struct {
			Operator string `json:"operator"`
		}
		if err := decodeJSON(r, &body); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		var item model.Alert
		var err error
		if parts[1] == "ack" {
			item, err = rt.lab.AcknowledgeAlert(r.Context(), parts[0], body.Operator)
		} else {
			item, err = rt.lab.ResolveAlert(r.Context(), parts[0], body.Operator)
		}
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	http.NotFound(w, r)
}
