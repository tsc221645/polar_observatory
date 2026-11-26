# Polar Observatory

Polar Observatory is a local-first Antarctic environmental data platform. The first vertical slice will follow daily Southern Hemisphere sea ice extent from the NSIDC source file through validation, PostgreSQL, the Go API, and the React explorer.

## Design baseline

The repository is at the start of v0.1.0. The following drafts define the first slice and its review gates.

- [Software specification](specs/spec_polar_observatory.md)
- [System context and component boundaries](docs/architecture/system.md)
- [Proposed architecture decision](docs/adr/0001-independent-react-and-grafana.md)
- [Source catalog](config/sources.yaml) and [data use review](DATA_LICENSES.md)
- [Draft OpenAPI contract](contracts/openapi/observations.yaml)
- [Draft sea ice explorer wireframe](docs/design/sea-ice-explorer.md)

The first bytes of the official NSIDC CSV have been verified. The v0.1.0 gate also requires product permissions, contract review, and an approved explorer mockup. The v0.2.0 milestone introduces the runnable React, Go, Python, and Compose baseline.

## Run the API locally

With Go 1.25 or later, run from the repository root:

```sh
cd services/api
go run ./cmd/server
```

In a second terminal, check liveness:

```sh
curl http://127.0.0.1:8000/api/v1/health/live
```

The response is `{"status":"ok"}`. Use `POLAR_API_ADDR` to override the default bind address. Run `go test ./...` and `go vet ./...` from `services/api`.

Once PostgreSQL is migrated and contains accepted observations, query the first
versioned data contract with an inclusive start and exclusive end date:

```sh
curl 'http://127.0.0.1:8000/api/v1/observations?metric=sea_ice_extent&entity=hemisphere%3Asouth&from=2024-01-01&to=2025-01-01&limit=100'
```

The response preserves the normalized value, quality flag, source and artifact
identifiers, and returns an opaque `next_cursor` when more rows are available.

## PostgreSQL readiness

`GET /api/v1/health/ready` returns HTTP 200 when the API can query the migrated source catalog. It returns HTTP 503 when PostgreSQL is unconfigured, unavailable, or the catalog is inaccessible. Liveness remains available during a database outage.

For a local database, use a PostgreSQL instance with PostGIS installed. From the repository root, set `GOOSE_DRIVER=postgres` and `GOOSE_DBSTRING` to the migration connection string and apply the initial migration:

```sh
goose -dir database/migrations up
```

The migration connection must be able to create the PostGIS extension and the `polar` schema. After applying it, grant an existing `api_reader` login access to the catalog in that database:

```sql
GRANT USAGE ON SCHEMA polar TO api_reader;
GRANT SELECT ON polar.sources TO api_reader;
```

Set `POLAR_API_DATABASE_URL` to the reader's connection string before starting the API. Check readiness with `curl http://127.0.0.1:8000/api/v1/health/ready`.
