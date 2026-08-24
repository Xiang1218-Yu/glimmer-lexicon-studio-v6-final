package api

import (
	"log/slog"
	"net/http"
	"time"

	"glimmer-lexicon-studio/internal/core"
)

type Server struct {
	engine *core.Engine
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(engine *core.Engine, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	server := &Server{engine: engine, logger: logger, mux: http.NewServeMux()}
	logger = logger.With("service", "glimmer-lexicon-studio")
	server.logger = logger
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler {
	return requestLog(s.logger, recoverPanic(s.mux))
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /readyz", s.ready)
	s.mux.HandleFunc("GET /v1/modules", s.modules)
	s.mux.HandleFunc("GET /v1/snapshot", s.snapshot)
	s.mux.HandleFunc("POST /v1/validate", s.validate)
	s.mux.HandleFunc("GET /v1/records", s.listRecords)
	s.mux.HandleFunc("POST /v1/records", s.createRecord)
	s.mux.HandleFunc("GET /v1/records/{id}", s.getRecord)
	s.mux.HandleFunc("POST /v1/records/{id}/advance", s.advanceRecord)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "glimmer-lexicon-studio", "time": time.Now().UTC()})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if s.engine == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"ready": false, "status": "unconfigured"})
		return
	}
	if !s.engine.Ready() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"ready": false, "status": "initializing"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ready": true})
}
