package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tsc221645/polar_observatory/services/api/internal/transport"
)

type Postgres struct{ pool *pgxpool.Pool }

func NewPostgres(pool *pgxpool.Pool) *Postgres { return &Postgres{pool: pool} }

func (s *Postgres) ListObservations(ctx context.Context, q transport.ObservationQuery) (transport.ObservationPage, error) {
	const statement = `
		WITH current_observations AS (
			SELECT o.*, row_number() OVER (
				PARTITION BY o.source_id, o.entity_key, o.metric, o.observed_at
				ORDER BY o.revision DESC, o.id DESC
			) AS revision_rank
			FROM polar.observations o
			WHERE o.entity_key = $1 AND o.metric = $2
			  AND o.observed_at >= $3 AND o.observed_at < $4
		)
		SELECT o.id, o.observed_at, o.value, o.unit, o.quality_flag,
		       o.source_id::text, o.artifact_id::text, s.version, o.ingested_at
		FROM current_observations o
		JOIN polar.sources s ON s.id = o.source_id
		WHERE o.revision_rank = 1
		  AND ($5::timestamptz IS NULL OR (o.observed_at, o.id) > ($5, $6))
		ORDER BY o.observed_at, o.id
		LIMIT $7`
	var afterTime any
	var afterID int64
	if q.After != nil {
		afterTime, afterID = q.After.Timestamp, q.After.ID
	}
	rows, err := s.pool.Query(ctx, statement, q.Entity, q.Metric, q.From, q.To, afterTime, afterID, q.Limit+1)
	if err != nil {
		return transport.ObservationPage{}, err
	}
	defer rows.Close()

	page := transport.ObservationPage{Data: make([]transport.Observation, 0, q.Limit), Meta: transport.PageMeta{}}
	var last transport.Cursor
	hasMore := false
	for rows.Next() {
		var id int64
		var item transport.Observation
		var version string
		var ingestedAt time.Time
		if err := rows.Scan(&id, &item.Timestamp, &item.Value, &item.Unit, &item.Quality, &item.SourceID, &item.ArtifactID, &version, &ingestedAt); err != nil {
			return transport.ObservationPage{}, err
		}
		page.Meta.SourceVersion = version
		if page.Meta.AsOf == nil || ingestedAt.After(*page.Meta.AsOf) {
			value := ingestedAt
			page.Meta.AsOf = &value
		}
		if len(page.Data) == q.Limit {
			hasMore = true
			break
		}
		page.Data = append(page.Data, item)
		last = transport.Cursor{Timestamp: item.Timestamp, ID: id}
	}
	if err := rows.Err(); err != nil {
		return transport.ObservationPage{}, err
	}
	if hasMore && last.ID != 0 {
		encoded := transport.EncodeCursor(last)
		page.NextCursor = &encoded
	}
	return page, nil
}

var _ transport.ObservationStore = (*Postgres)(nil)
