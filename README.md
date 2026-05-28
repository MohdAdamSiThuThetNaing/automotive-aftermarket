# Automotive Aftermarket Translation Loader

Take-home backend exercise for WYZauto.

This project implements a production-oriented translation loading layer in Go for assembling Elasticsearch product documents from PostgreSQL translation data.

---

# Features

- Bulk translation loading (no N+1 queries)
- Locale filtering support (`en`, `th`, `ko`)
- Thread-safe in-memory TTL cache
- Entity-level cache invalidation
- Graceful translation fallback (`en` → empty string)
- Dockerized PostgreSQL setup
- Unit tests using interfaces and mocks
- Integration test using real PostgreSQL
- Clean layered architecture

---

# Tech Stack

- Go
- PostgreSQL
- pgx/v5
- Docker Compose

Additional dependency:

- `github.com/joho/godotenv`

Used only for local environment variable loading during development.
Core translation loading logic depends only on `pgx`.

---

# Project Structure

```txt
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── app
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── builder
│   │   └── product_document_builder.go
│   ├── cache
│   │   └── translation_cache.go
│   ├── config
│   │   └── config.go
│   ├── loader
│   │   ├── postgres_loader.go
│   │   └── translation_loader.go
│   ├── mocks
│   │   ├── mock_product_repository.go
│   │   └── mock_translation_loader.go
│   ├── models
│   │   ├── elasticsearch_document.go
│   │   ├── product.go
│   │   ├── specification.go
│   │   └── translation.go
│   └── repository
│       ├── postgres_product_repository.go
│       └── product_repository.go
├── migrations
│   ├── 1-init.sql
│   └── 2-seed.sql
└── tests
    ├── integration_test.go
    └── product_document_builder_test.go
```

---

# Architecture Overview

The project is separated into dedicated layers:

| Layer      | Responsibility                  |
| ---------- | ------------------------------- |
| repository | Database access                 |
| loader     | Bulk translation loading        |
| cache      | In-memory TTL caching           |
| builder    | Elasticsearch document assembly |
| mocks      | Mock implementations for tests  |
| tests      | Unit + integration tests        |

This separation keeps the translation loader testable and isolates database concerns behind interfaces.

---

# Query Strategy

The translation loader avoids N+1 queries by bulk-loading translations using a single SQL query:

```sql
SELECT
    entity_type,
    entity_id,
    locale,
    field_name,
    field_value,
    updated_at
FROM translation
WHERE entity_type = $1
AND entity_id = ANY($2)
AND locale = ANY($3)
```

This allows translations for multiple entities and locales to be loaded in a single database round-trip.

---

# Cache Design

The translation loader includes a thread-safe in-memory cache.

Features:

- TTL-based expiration
- `sync.RWMutex` for concurrent access safety
- Entity-level invalidation
- Lazy expiration strategy

Example cache key:

```txt
product:11111111-1111-1111-1111-111111111111
```

Trade-off:

- Short-lived stale reads are possible during concurrent sync operations.
- This is acceptable because Elasticsearch indexing is eventually consistent.
- Short TTL minimizes stale exposure window while avoiding complex distributed invalidation mechanisms.

---

# Graceful Fallback Strategy

Missing translations never return an error or panic.

Fallback order:

1. Requested locale
2. English (`en`)
3. Empty string

Example:

- Thai translation missing
- English translation returned automatically

---

# How To Run

## 1. Clone Repository

```bash
git clone https://github.com/MohdAdamSiThuThetNaing/automotive-aftermarket.git
cd automotive-aftermarket
```

---

## 2. Start PostgreSQL + Application

```bash
docker compose up --build
```

This automatically:

- starts PostgreSQL
- initializes schema
- seeds test data
- runs the application

---

# Environment Variables

Example `.env`:

```env
DATABASE_URL=

PRODUCT_ID=11111111-1111-1111-1111-111111111111

APP_ENV=development
APP_NAME=automotive-aftermarket

CACHE_TTL_MINUTES=5
```

---

# Database Initialization

Schema and seed data are automatically initialized using Docker entrypoint scripts:

```txt
migrations/1-init.sql
migrations/2-seed.sql
```

---

# Running Tests

## Run All Tests

```bash
go test ./... -v
```

---

## Unit Tests

Unit tests:

- use interfaces and mocks
- require no database
- test business logic in isolation

Example:

- fallback to English translation

---

## Integration Test

Integration test:

- uses real PostgreSQL
- uses seeded data
- validates full document assembly flow

---

# Example Output

```json
{
  "uuid": "11111111-1111-1111-1111-111111111111",
  "sku": "BP-OIL-5W30-1L",
  "part_number": "5W30-1L",
  "productname": [
    {
      "locale": "en",
      "data": "5W-30 Engine Oil 1L"
    },
    {
      "locale": "th",
      "data": "น้ำมันเครื่อง 5W-30 1 ลิตร"
    }
  ]
}
```

---

# Design Question — Delta Sync Strategy

To support delta sync, I would extend the loader with a query filtered by `updated_at`.

Example strategy:

```sql
SELECT
    entity_type,
    entity_id,
    locale,
    field_name,
    field_value,
    updated_at
FROM translation
WHERE updated_at > $1
ORDER BY updated_at ASC
LIMIT $2
```

The sync pipeline would persist the latest processed timestamp as a cursor.
Future syncs would only reload translations updated after that timestamp.

Main trade-off:

- timestamp-based polling is simpler and lightweight
- but introduces eventual consistency concerns
- updates occurring during sync execution may appear in the next sync cycle

For larger scale systems, I would eventually move toward a Change Data Capture (CDC) approach using PostgreSQL WAL streaming or Debezium.

---

# Future Improvements

Given more time, I would add:

- Redis distributed cache
- Metrics and tracing
- Structured logging
- Background cache refresh
- Prepared statement optimization
- CDC / Debezium integration
- Elasticsearch bulk indexing pipeline
- Benchmark tests
- Generic translation aggregation utilities

---

# Notes

This implementation prioritizes:

- clarity
- testability
- operational simplicity
- maintainability

over excessive abstraction or framework complexity.
