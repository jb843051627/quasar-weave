package handler

import (
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
)

func (rt *Router) handleObservations(w http.ResponseWriter, r *http.Request, suffix string) {
	parts := splitPath(suffix)
	if len(parts) == 0 {
		if r.Method == http.MethodGet {
			filter := model.ObservationFilter{Target: r.URL.Query().Get("target"), Status: model.ObservationStatus(r.URL.Query().Get("status"))}
			items, err := rt.lab.ListObservations(r.Context(), filter)
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
		var input model.ObservationInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := rt.lab.CreateObservation(r.Context(), input)
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
		return
	}
	id := parts[0]
	if len(parts) == 1 {
		if !method(w, r, http.MethodGet) {
			return
		}
		item, err := rt.lab.GetObservation(r.Context(), id)
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	if len(parts) == 2 {
		if parts[1] == "summary" {
			if !method(w, r, http.MethodGet) {
				return
			}
			item, err := rt.lab.ObservationSummary(r.Context(), id)
			if err != nil {
				writeError(w, service.StatusCode(err), err)
				return
			}
			writeOK(w, item)
			return
		}
		if parts[1] == "frames" {
			if !method(w, r, http.MethodGet) {
				return
			}
			items, err := rt.lab.ListFrames(r.Context(), model.FrameFilter{ObservationID: id})
			if err != nil {
				writeError(w, service.StatusCode(err), err)
				return
			}
			writeOK(w, items)
			return
		}
		var item model.Observation
		var err error
		switch parts[1] {
		case "start":
			if !method(w, r, http.MethodPost) {
				return
			}
			item, err = rt.lab.StartObservation(r.Context(), id)
		case "calibrate":
			if !method(w, r, http.MethodPost) {
				return
			}
			item, err = rt.lab.BeginCalibration(r.Context(), id)
		case "archive":
			if !method(w, r, http.MethodPost) {
				return
			}
			item, err = rt.lab.ArchiveObservation(r.Context(), id)
		default:
			http.NotFound(w, r)
			return
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
