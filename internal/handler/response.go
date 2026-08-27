package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/service"
)

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, args ...any) {
	status := http.StatusInternalServerError
	var err error
	if len(args) == 1 {
		err, _ = args[0].(error)
	} else if len(args) >= 2 {
		status, _ = args[0].(int)
		err, _ = args[1].(error)
	}
	if err == nil {
		err = fmt.Errorf("request failed")
	}
	if len(args) == 1 {
		status = service.StatusCode(err)
	}
	writeJSON(w, status, errorBody{Error: err.Error()})
}

func writeNoContent(w http.ResponseWriter) { w.WriteHeader(http.StatusNoContent) }

func writeOK(w http.ResponseWriter, value any) { writeJSON(w, http.StatusOK, value) }
