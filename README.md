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
