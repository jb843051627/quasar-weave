package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func decodeJSON(r *http.Request, target any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	return nil
}

func decodeBody(r *http.Request, target any) error { return decodeJSON(r, target) }

func method(w http.ResponseWriter, r *http.Request, expected ...string) bool {
	for _, candidate := range expected {
		if r.Method == candidate {
			return true
		}
	}
	allow := ""
	for i, candidate := range expected {
		if i > 0 {
			allow += ", "
		}
		allow += candidate
	}
	w.Header().Set("Allow", allow)
	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}
