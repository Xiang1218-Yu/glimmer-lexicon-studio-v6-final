package api

import (
	"net/http"
)

func (s *Server) modules(w http.ResponseWriter, r *http.Request) {
	modules := s.engine.Modules()
	values := make([]map[string]any, 0, len(modules))
	for _, module := range modules {
		if !module.Enabled() {
			continue
		}
		values = append(values, map[string]any{
			"key": module.Key(), "description": module.Description(),
			"priority": module.Priority(), "family": module.Family(), "enabled": module.Enabled(),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": values, "count": len(values)})
}

func (s *Server) snapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.engine.Snapshot())
}
