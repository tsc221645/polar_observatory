package transport

import (
	"context"
	"net/http"
	"time"
)

func NewHandler(checkReady func(context.Context) bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health/live", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	mux.HandleFunc("GET /api/v1/health/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if checkReady == nil || !checkReady(ctx) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("{\"status\":\"unavailable\"}\n"))
			return
		}
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	return mux
}
