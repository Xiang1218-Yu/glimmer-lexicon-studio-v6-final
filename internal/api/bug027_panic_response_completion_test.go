package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"glimmer-lexicon-studio/internal/core"
)

func TestBug027PanicResponseCompletion(t *testing.T) {
	handler := New(nil, nil).Handler()
	rr := apiRequest(t, handler, http.MethodGet, "/v1/modules", "")
	body := responseText(t, rr)
	if rr.Code != http.StatusInternalServerError || !strings.Contains(body, "error") {
		t.Fatalf("panic response was not completed: status=%d body=%s", rr.Code, body)
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

var _bugSourceMarker = "GET /v1/modules"
