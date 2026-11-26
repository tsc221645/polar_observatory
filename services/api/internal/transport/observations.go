package transport

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

const maxObservationLimit = 1000

type Observation struct {
	Timestamp  time.Time `json:"timestamp"`
	Value      *float64  `json:"value"`
	Unit       string    `json:"unit"`
	Quality    string    `json:"quality"`
	SourceID   string    `json:"source_id"`
	ArtifactID string    `json:"artifact_id"`
}

type ObservationQuery struct {
	Metric string
	Entity string
	From   time.Time
	To     time.Time
	Limit  int
	After  *Cursor
}

type Cursor struct {
	Timestamp time.Time
	ID        int64
}

type ObservationPage struct {
	Data       []Observation `json:"data"`
	NextCursor *string       `json:"next_cursor"`
	Meta       PageMeta      `json:"meta"`
}

type PageMeta struct {
	SourceVersion string     `json:"source_version"`
	AsOf          *time.Time `json:"as_of"`
}

type ObservationStore interface {
	ListObservations(context.Context, ObservationQuery) (ObservationPage, error)
}

func observationHandler(store ObservationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := parseObservationQuery(r)
		if err != nil {
			writeAPIError(w, http.StatusBadRequest, "invalid_query", err.Error())
			return
		}
		if query.Metric != "sea_ice_extent" {
			writeAPIError(w, http.StatusUnprocessableEntity, "unsupported_metric", "metric must be sea_ice_extent")
			return
		}
		if query.To.Before(query.From) || query.To.Equal(query.From) || query.To.After(query.From.AddDate(10, 0, 0)) {
			writeAPIError(w, http.StatusUnprocessableEntity, "invalid_range", "date range must be positive and at most ten years")
			return
		}
		if store == nil {
			writeAPIError(w, http.StatusServiceUnavailable, "database_unavailable", "observation store is unavailable")
			return
		}
		page, err := store.ListObservations(r.Context(), query)
		if err != nil {
			writeAPIError(w, http.StatusServiceUnavailable, "database_unavailable", "observations are temporarily unavailable")
			return
		}
		writeJSON(w, http.StatusOK, page)
	}
}

func parseObservationQuery(r *http.Request) (ObservationQuery, error) {
	values := r.URL.Query()
	q := ObservationQuery{Metric: values.Get("metric"), Entity: values.Get("entity"), Limit: 1000}
	if q.Metric == "" || q.Entity == "" {
		return q, errors.New("metric and entity are required")
	}
	from, err := time.Parse("2006-01-02", values.Get("from"))
	if err != nil {
		return q, errors.New("from must be an ISO-8601 date")
	}
	to, err := time.Parse("2006-01-02", values.Get("to"))
	if err != nil {
		return q, errors.New("to must be an ISO-8601 date")
	}
	q.From, q.To = from.UTC(), to.UTC()
	if raw := values.Get("limit"); raw != "" {
		q.Limit, err = strconv.Atoi(raw)
		if err != nil || q.Limit < 1 || q.Limit > maxObservationLimit {
			return q, errors.New("limit must be between 1 and 1000")
		}
	}
	if raw := values.Get("cursor"); raw != "" {
		q.After, err = decodeCursor(raw)
		if err != nil {
			return q, errors.New("cursor is invalid")
		}
	}
	return q, nil
}

func encodeCursor(cursor Cursor) string {
	var payload [16]byte
	binary.BigEndian.PutUint64(payload[:8], uint64(cursor.Timestamp.UnixNano()))
	binary.BigEndian.PutUint64(payload[8:], uint64(cursor.ID))
	mac := hmac.New(sha256.New, []byte("polar-observatory-cursor-v1"))
	mac.Write(payload[:])
	return base64.RawURLEncoding.EncodeToString(append(payload[:], mac.Sum(nil)[:8]...))
}

// EncodeCursor is used by repository implementations to create the opaque
// continuation token returned by the HTTP contract.
func EncodeCursor(cursor Cursor) string { return encodeCursor(cursor) }

func decodeCursor(raw string) (*Cursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil || len(b) != 24 {
		return nil, errors.New("bad cursor")
	}
	mac := hmac.New(sha256.New, []byte("polar-observatory-cursor-v1"))
	mac.Write(b[:16])
	if !hmac.Equal(b[16:], mac.Sum(nil)[:8]) {
		return nil, errors.New("bad cursor")
	}
	return &Cursor{Timestamp: time.Unix(0, int64(binary.BigEndian.Uint64(b[:8]))).UTC(), ID: int64(binary.BigEndian.Uint64(b[8:16]))}, nil
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message, "request_id": ""}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
