package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"glimmer-lexicon-studio/internal/core"
)

func TestBug026RequestBodyClosedOnDecodeError(t *testing.T) {
	handler := New(core.NewEngine(), nil).Handler()
	req := httptest.NewRequest(http.MethodPost, "/v1/records", strings.NewReader("{"))
	body := &trackingBody{Reader: req.Body}
	req.Body = body
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rr.Code)
	}
	if !body.closed {
		t.Fatal("request body was not closed on the decoder error path")
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
