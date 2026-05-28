# Automotive Aftermarket Translation Loader

Take-home backend exercise for WYZauto.

This project implements a production-oriented translation loading layer in Go for assembling Elasticsearch product documents from PostgreSQL translation data.

---

# Features

- Bulk translation loading (no N+1 queries)
- Locale filtering support (`en`, `th`, `ko`)
- Thread-safe in-process in-memory TTL cache
- Entity-level cache invalidation
- Graceful translation fallback (`en` → empty string)
- Product specification aggregation
- Attribute translation support
- Elasticsearch-ready document assembly
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
- `golang.org/x/sync/errgroup`

Used only for local environment variable loading and lightweight concurrency orchestration.

Core translation loading logic depends only on `pgx`.

---

# Project Structure

```txt id="tnw5ux"
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── app
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── builder
│   │   └── product_document_builder.go
│   ├── cache
│   │   └── translation_cache.go
│   ├── config
│   │   └── config.go
│   ├── loader
│   │   ├── postgres_loader.go
│   │   └── translation_loader.go
│   ├── mocks
│   │   ├── mock_product_repository.go
│   │   └── mock_translation_loader.go
│   ├── models
│   │   ├── elasticsearch_document.go
│   │   ├── product.go
│   │   ├── specification.go
│   │   └── translation.go
│   └── repository
│       ├── postgres_product_repository.go
│       └── product_repository.go
├── migrations
│   ├── 1-init.sql
│   └── 2-seed.sql
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
| cache      | In-process TTL caching          |
| builder    | Elasticsearch document assembly |
| mocks      | Mock implementations for tests  |
| tests      | Unit + integration tests        |

This separation keeps the translation loader testable while isolating persistence concerns behind interfaces.

---

# Design Decisions

Several design decisions were intentionally made to align with the exercise requirements while keeping the implementation maintainable and operationally simple.

---

## Why Interfaces?

The repository and translation loader are abstracted behind interfaces to:

- isolate persistence concerns
- improve testability
- support mock-based unit testing
- allow future backend replacement with minimal business-logic changes

---

## Why Bulk Translation Loading?

The primary performance concern described in the exercise was N+1 translation queries.

The loader therefore batches translation retrieval using:

```sql id="scv3dj"
WHERE entity_id = ANY($2)
```

This minimizes database round-trips and scales significantly better during Elasticsearch indexing operations.

---

## Why In-Process Cache Instead of Redis?

The exercise explicitly requested:

- thread-safe in-process cache
- TTL support
- entity-level invalidation

An in-memory cache provides:

- low latency
- operational simplicity
- zero external infrastructure requirements

Redis would be more appropriate for:

- distributed workers
- multi-instance synchronization
- shared cache consistency

but would add operational complexity outside the scope of this exercise.

---

## Why Graceful Translation Fallback?

Search indexing pipelines should remain resilient even when translation coverage is incomplete.

The implementation therefore prioritizes:

- successful document generation
- graceful degradation
- operational continuity

instead of strict translation completeness.

---

## Why errgroup?

`errgroup` allows independent I/O operations to execute concurrently while maintaining:

- clean cancellation handling
- structured error propagation
- simpler concurrency management

This improves overall document assembly latency without introducing unnecessary complexity.

---

# Translation Loading Strategy

The translation loader avoids N+1 queries by bulk-loading translations using a single database round-trip.

Example query:

```sql id="06q4bd"
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

This allows translations for multiple entities and locales to be loaded efficiently in a single query.

---

# Elasticsearch Document Assembly

The builder layer assembles a denormalized Elasticsearch-ready document from:

- product data
- translations
- product specifications
- attribute metadata

Example output:

```json id="3xhnx6"
{
  "uuid": "11111111-1111-1111-1111-111111111111",
  "sku": "BP-OIL-5W30-1L",
  "part_number": "5W30-1L",
  "brand": {
    "code": "bosch",
    "label": {}
  },
  "productname": [
    {
      "locale": "en",
      "data": "5W-30 Engine Oil 1L"
    },
    {
      "locale": "th",
      "data": "น้ำมันเครื่อง 5W-30 1 ลิตร"
    }
  ],
  "attributes": {
    "oil_grade": "5w30"
  },
  "oil_grade": {
    "code": "5w30",
    "label": {
      "en": "5W-30",
      "th": "5W-30"
    }
  }
}
```

---

# Cache Design

The translation loader uses a thread-safe in-process in-memory TTL cache.

The exercise explicitly requested an in-process caching layer, so the implementation intentionally avoids external distributed cache systems such as Redis in favor of a lightweight and operationally simple design.

Features:

- Go map-based storage
- `sync.RWMutex` for concurrent access safety
- TTL-based expiration
- Entity-level invalidation
- Lazy expiration strategy

Example cache key:

```txt id="o2uv2u"
product:11111111-1111-1111-1111-111111111111
```

Entity-level invalidation:

```go id="8nk7hu"
Invalidate(entityType, entityID)
```

Trade-offs:

- short-lived stale reads are possible while cache entries remain valid
- cache is scoped per application instance
- cache is cleared on application restart

These trade-offs are acceptable because Elasticsearch synchronization pipelines are typically eventually consistent.

The design prioritizes:

- simplicity
- maintainability
- low operational overhead
- fast in-process reads
- clear invalidation semantics

---

# Concurrent Document Loading

`ProductDocumentBuilder` uses `errgroup` to parallelize independent I/O operations:

- translation loading
- specification loading

This reduces total document assembly latency while keeping error propagation clean and structured.

---

# Attribute Document Strategy

The implementation currently supports structured attribute rendering for:

- `oil_grade`

using:

```json id="2qfqot"
{
  "oil_grade": {
    "code": "5w30",
    "label": {
      "en": "5W-30",
      "th": "5W-30"
    }
  }
}
```

Assumption:

The exercise document only provided one explicit example of structured specification rendering (`oil_grade`).

In a production system, this would likely evolve into:

- schema-driven attribute rendering
- configuration-based attribute mapping
- dynamic attribute serializers

instead of hardcoded switch logic.

---

# Graceful Fallback Strategy

Missing translations never return an error or panic.

Fallback order:

1. Requested locale
2. English (`en`)
3. Empty string

Example:

- Thai translation missing
- English translation automatically returned

---

# Error Handling Strategy

Errors are wrapped with contextual information:

```go id="4zjlwm"
return nil, fmt.Errorf(
    "load product translations: %w",
    err,
)
```

This improves:

- observability
- debugging
- operational diagnostics

while preserving original error chains.

---

# Environment Variables

Example `.env`:

```env id="ry9e9p"
DATABASE_URL=

PRODUCT_ID=11111111-1111-1111-1111-111111111111

APP_ENV=development
APP_NAME=automotive-aftermarket

CACHE_TTL_MINUTES=5
```

---

# Database Initialization

Schema and seed data are automatically initialized using Docker entrypoint scripts:

```txt id="0e8d77"
migrations/1-init.sql
migrations/2-seed.sql
```

---

# Running The Project

## Start PostgreSQL + Application

```bash id="4eyg1l"
make up
```

OR:

```bash id="c4ny2m"
docker compose up --build
```

This automatically:

- starts PostgreSQL
- initializes schema
- seeds sample data
- runs the application

---

# Running Tests

## Run All Tests

```bash id="8pf7y9"
make test
```

OR:

```bash id="4kr5rq"
go test ./... -v
```

---

# Unit Tests

Unit tests:

- use interfaces and mocks
- require no database
- test business logic in isolation

Example:

- graceful English fallback behavior

---

# Integration Test

Integration test:

- uses real PostgreSQL
- uses seeded data
- validates full document assembly flow
- validates translation loading behavior

---

# Makefile Commands

```bash id="kl6gln"
make up
```

Start Docker services.

```bash id="43ck0f"
make down
```

Stop Docker services.

```bash id="f1y70g"
make run
```

Run application locally.

```bash id="m1f8nq"
make test
```

Run all tests.

---

# Design Question — Delta Sync Strategy

To support delta sync, I would extend the loader using an `updated_at` cursor strategy.

Example query:

```sql id="y00ls3"
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

The synchronization pipeline would persist the latest processed timestamp as a cursor.

Future synchronization runs would only reload translations updated after the stored timestamp.

Trade-offs:

- timestamp-based polling is simple and operationally lightweight
- introduces eventual consistency behavior
- updates occurring during synchronization may appear in the next sync cycle

For larger-scale systems, I would eventually move toward a CDC-based approach using PostgreSQL WAL streaming or Debezium.

---

# Future Improvements

Given more time, I would add:

- Redis distributed cache
- Metrics and tracing
- Structured logging
- Background cache refresh
- Prepared statement optimization
- Generic translation aggregation utilities
- Elasticsearch bulk indexing pipeline
- CDC / Debezium integration
- Benchmark tests
- OpenTelemetry support

---

# Assumptions

This implementation makes several assumptions:

- Elasticsearch indexing is eventually consistent
- cache is scoped to a single application instance
- translation volume is manageable within batch-loading strategy
- one product document is assembled per request
- attribute rendering requirements are limited to provided examples

These assumptions keep the implementation intentionally focused and aligned with the exercise scope.

---

# Notes

This implementation prioritizes:

- maintainability
- operational simplicity
- testability
- clear abstractions
- production-oriented query behavior

over excessive framework complexity or premature optimization.
