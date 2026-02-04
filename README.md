# Pack Calculator API (Go)

HTTP API + simple UI to compute optimal pack combinations for a given order amount.

## Overview
- Primary objective: minimize items shipped (overage).
- Secondary objective: minimize number of packs.
- Pack sizes are configurable at runtime via API.


## Quickstart (Local)
```bash
go run ./cmd/server
```

Open:
- UI: http://localhost:8080
- Health: http://localhost:8080/healthz

## Docker Compose
```bash
docker compose up --build
```

## API
- `GET /healthz`
- `GET /v1/pack-sizes`
- `PUT /v1/pack-sizes`
- `POST /v1/calculate`

Example:
```bash
curl -s -X PUT http://localhost:8080/v1/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{"pack_sizes":[23,31,53]}'

curl -s -X POST http://localhost:8080/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"amount":500000}'
```

## UI
Open `http://localhost:8080` in the browser.

## Tests
```bash
go test ./...
```

## Configuration
- `PORT` (default: `8080`)
- `DB_PATH` (default: `data/packs.db`)
- `WEB_DIR` (default: `web`)
