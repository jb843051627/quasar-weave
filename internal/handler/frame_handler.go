package handler

import (
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/model"
	"github.com/jb843051627/quasar-weave/internal/service"
)

func (rt *Router) handleFrames(w http.ResponseWriter, r *http.Request, suffix string) {
	parts := splitPath(suffix)
	if len(parts) == 0 {
		if r.Method == http.MethodGet {
			filter := model.FrameFilter{ObservationID: r.URL.Query().Get("observation_id"), AntennaID: r.URL.Query().Get("antenna_id"), State: model.FrameState(r.URL.Query().Get("state"))}
			items, err := rt.lab.ListFrames(r.Context(), filter)
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
		var input model.FrameInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, done, err := rt.lab.SubmitFrame(r.Context(), input)
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		if done != nil {
			select {
			case processingErr := <-done:
				if processingErr != nil {
					writeError(w, service.StatusCode(processingErr), processingErr)
					return
				}
			case <-r.Context().Done():
				writeError(w, http.StatusRequestTimeout, r.Context().Err())
				return
			}
		}
		writeJSON(w, http.StatusAccepted, item)
		return
	}
	if len(parts) == 1 {
		if !method(w, r, http.MethodGet) {
			return
		}
		item, err := rt.lab.GetFrame(r.Context(), parts[0])
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	if len(parts) == 2 && parts[1] == "evaluate" {
		if !method(w, r, http.MethodPost) {
			return
		}
		item, err := rt.lab.EvaluateFrame(r.Context(), parts[0])
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeOK(w, item)
		return
	}
	http.NotFound(w, r)
}
