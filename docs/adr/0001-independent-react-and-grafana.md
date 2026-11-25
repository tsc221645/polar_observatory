# 0001. Independent React and Grafana interfaces

**Date** 2026-11-25  
**Status** Proposed

## Context and drivers

Public visitors need accessible exploration, shareable filters, and provenance alongside observations. Scientists and operators need flexible dashboards and telemetry. The public explorer must stay available during Grafana maintenance.

## Options and proposed decision

Use a React public explorer backed by the Go API and provision Grafana separately for scientific and operational views. A single Grafana interface would reduce UI work but give the public explorer less control over navigation, accessibility, and provenance. Embedding Grafana panels in React would introduce a runtime dependency between the two interfaces.

## Consequences and evidence

Both interfaces need deployment and smoke coverage. Shared scientific definitions live in the publication layer and API contract. Verify the first public explorer against a synthetic fixture with Grafana stopped, then verify Grafana independently against its read-only data source when that milestone lands.

Revisit the decision if independent operation fails in the Compose smoke test or if the public explorer cannot meet the accessibility acceptance criteria in the [specification](../../specs/spec_polar_observatory.md).
