package api

import (
	"context"
	"errors"
	"net/http"

	"glimmer-lexicon-studio/internal/core"
)

type recordInput struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

type advanceInput struct {
	Stage string `json:"stage"`
}

func (s *Server) listRecords(w http.ResponseWriter, r *http.Request) {
	records, err := s.engine.List(r.Context(), r.URL.Query().Get("stage"), limit(r.URL.Query().Get("limit"), 100))
	if err != nil {
		// A cancelled/deadline-exceeded context means the scan could not finish
		// in time. Surface it as a transient failure rather than writing back a
		// partial, stale list that would mislead the next page refresh.
		if r.Context().Err() != nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			writeError(w, http.StatusServiceUnavailable, "list scan was cancelled, please retry")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": records, "count": len(records)})
}

func (s *Server) createRecord(w http.ResponseWriter, r *http.Request) {
	var input recordInput
	if !decode(w, r, &input) {
		return
	}
	record, err := s.engine.Create(r.Context(), input.ID, input.Payload, actor(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) getRecord(w http.ResponseWriter, r *http.Request) {
	record, err := s.engine.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) advanceRecord(w http.ResponseWriter, r *http.Request) {
	var input advanceInput
	if !decode(w, r, &input) {
		return
	}
	record, err := s.engine.Advance(r.Context(), r.PathValue("id"), input.Stage, actor(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func _keepCoreType(_ core.Record) {}
