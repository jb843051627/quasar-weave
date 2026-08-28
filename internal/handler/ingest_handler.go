package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/jb843051627/quasar-weave/internal/protocol"
	"github.com/jb843051627/quasar-weave/internal/service"
)

func (h *Router) handleIngest(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodPost) {
		return
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-quasar-lines") {
		source := r.URL.Query().Get("source")
		result, err := h.lab.IngestStream(r.Context(), r.Body, source, time.Now().UTC())
		if err != nil {
			writeError(w, service.StatusCode(err), err)
			return
		}
		writeJSON(w, http.StatusAccepted, result)
		return
	}
	var envelope protocol.Envelope
	if err := decodeJSON(r, &envelope); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.lab.IngestEnvelope(r.Context(), envelope)
	if err != nil {
		writeError(w, service.StatusCode(err), err)
		return
	}
	writeJSON(w, http.StatusAccepted, result)
}
