package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"glimmer-lexicon-studio/internal/core"
)

func TestBug030CreateErrorStatus(t *testing.T) {
	handler := New(core.NewEngine(), nil).Handler()
	first := apiRequest(t, handler, http.MethodPost, "/v1/records", `{"id":"duplicate","payload":"first"}`)
	if first.Code != http.StatusCreated {
		t.Fatalf("first create failed: %d %s", first.Code, responseText(t, first))
	}
	second := apiRequest(t, handler, http.MethodPost, "/v1/records", `{"id":"duplicate","payload":"second"}`)
	if second.Code != http.StatusBadRequest {
		t.Fatalf("expected client error for duplicate record, got %d body=%s", second.Code, responseText(t, second))
	}
}

func apiRequest(t *testing.T, handler http.Handler, method, path, payload string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(payload))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func responseText(t *testing.T, rr *httptest.ResponseRecorder) string {
	t.Helper()
	data, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	return string(data)
}

type trackingBody struct {
	io.Reader
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

var _ = core.NewEngine

var _bugSourceMarker = "POST /v1/records"
