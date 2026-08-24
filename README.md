# Glimmer Lexicon Studio

Glimmer coordinates concepts, multilingual terms, contextual usage, language review, replacement relations and immutable release lexicons.

This is a standalone Go service using only the standard library. It keeps a
small in-memory domain store so the complete lifecycle can be explored without
external infrastructure.

## Run

```bash
go run ./cmd/glimmer-lexicon-studio
```

The default address is `:8282`; set `SERVICE_ADDR` to override it.

## Quick flow

```bash
curl http://localhost:8282/healthz
curl -X POST http://localhost:8282/v1/records \
  -H 'content-type: application/json' \
  -d '{"id":"demo-1","payload":"concept payment settlement preferred term"}'
curl -X POST http://localhost:8282/v1/records/demo-1/advance \
  -H 'content-type: application/json' \
  -d '{"stage":"review"}'
```

The project has no test files by design; package compilation, vetting and
runtime API checks are used for validation.
