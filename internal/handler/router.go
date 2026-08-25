package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jb843051627/quasar-weave/internal/service"
)

type Router struct {
	lab *service.Lab
}

func New(lab *service.Lab) http.Handler { return &Router{lab: lab} }

func (h *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	defer func() { log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(started)) }()
	w.Header().Set("X-Content-Type-Options", "nosniff")
	path := strings.TrimSuffix(r.URL.Path, "/")
	switch {
	case path == "" || path == "/index.html":
		h.handleWeb(w, r)
	case path == "/healthz":
		h.healthz(w, r)
	case path == "/readyz":
		h.ready(w, r)
	case path == "/api/health":
		h.health(w, r)
	case strings.HasPrefix(path, "/api/antennas"):
		h.antennas(w, r, strings.TrimPrefix(path, "/api/antennas"))
	case strings.HasPrefix(path, "/api/gates"):
		h.handleGates(w, r, strings.TrimPrefix(path, "/api/gates"))
	case strings.HasPrefix(path, "/api/observations"):
		h.handleObservations(w, r, strings.TrimPrefix(path, "/api/observations"))
	case strings.HasPrefix(path, "/api/frames"):
		h.handleFrames(w, r, strings.TrimPrefix(path, "/api/frames"))
	case strings.HasPrefix(path, "/api/alerts"):
		h.handleAlerts(w, r, strings.TrimPrefix(path, "/api/alerts"))
	case strings.HasPrefix(path, "/api/notes"):
		h.handleNotes(w, r, strings.TrimPrefix(path, "/api/notes"))
	case strings.HasPrefix(path, "/api/events"):
		h.handleEvents(w, r)
	default:
		http.NotFound(w, r)
	}
}
