package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type fakeObservationStore struct {
	page  ObservationPage
	query ObservationQuery
	err   error
}

func (f *fakeObservationStore) ListObservations(_ context.Context, query ObservationQuery) (ObservationPage, error) {
	f.query = query
	return f.page, f.err
}

func TestObservationsQueryValidationAndStoreCall(t *testing.T) {
	store := &fakeObservationStore{page: ObservationPage{Data: []Observation{}}}
	handler := observationHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/observations?metric=sea_ice_extent&entity=hemisphere%3Asouth&from=2024-01-01&to=2024-02-01&limit=25", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if store.query.Limit != 25 || store.query.Entity != "hemisphere:south" {
		t.Fatalf("query = %+v", store.query)
	}
}

func TestObservationsErrorsUseContract(t *testing.T) {
	for _, test := range []struct {
		name, target string
		status       int
	}{
		{"missing required dates", "/api/v1/observations?metric=sea_ice_extent&entity=x", http.StatusBadRequest},
		{"unsupported metric", "/api/v1/observations?metric=temperature&entity=x&from=2024-01-01&to=2024-02-01", http.StatusUnprocessableEntity},
		{"invalid cursor", "/api/v1/observations?metric=sea_ice_extent&entity=x&from=2024-01-01&to=2024-02-01&cursor=bad", http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			observationHandler(&fakeObservationStore{}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.target, nil))
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
			var body map[string]map[string]string
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil || body["error"]["code"] == "" {
				t.Fatalf("error body = %s", response.Body.String())
			}
		})
	}
}

func TestCursorIsOpaqueAndSigned(t *testing.T) {
	encoded := EncodeCursor(Cursor{Timestamp: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), ID: 42})
	decoded, err := decodeCursor(encoded)
	if err != nil || decoded.ID != 42 {
		t.Fatalf("decoded cursor = %+v, err = %v", decoded, err)
	}
	if _, err := decodeCursor(strings.TrimSuffix(encoded, "A") + "B"); err == nil {
		t.Fatal("tampered cursor was accepted")
	}
}
