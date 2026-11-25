package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveness(t *testing.T) {
	handler := NewHandler(nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/health/live", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := response.Body.String(); got != "{\"status\":\"ok\"}\n" {
		t.Errorf("body = %q, want liveness status", got)
	}
}

func TestReadiness(t *testing.T) {
	for _, test := range []struct {
		name       string
		checkReady func(context.Context) bool
		wantStatus int
		wantBody   string
	}{
		{"unconfigured", nil, http.StatusServiceUnavailable, "{\"status\":\"unavailable\"}\n"},
		{"database check failed", func(context.Context) bool { return false }, http.StatusServiceUnavailable, "{\"status\":\"unavailable\"}\n"},
		{"database ready", func(context.Context) bool { return true }, http.StatusOK, "{\"status\":\"ok\"}\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
			NewHandler(test.checkReady).ServeHTTP(response, request)
			if response.Code != test.wantStatus || response.Body.String() != test.wantBody {
				t.Fatalf("response = %d %q, want %d %q", response.Code, response.Body.String(), test.wantStatus, test.wantBody)
			}
		})
	}
}
