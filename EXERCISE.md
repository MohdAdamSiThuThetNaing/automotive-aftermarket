# WYZauto Backend Take-Home Exercise

## Context

WYZauto is a B2B automotive aftermarket marketplace operating across Thailand, Malaysia, and South Korea. Product catalogue data is stored in PostgreSQL, while buyer-facing search is powered by Elasticsearch.

Product content is multi-locale (`en`, `th`, `ko`). Product names, descriptions, attribute labels, and specification value labels are stored as rows in a generic `translation` table.

During indexing, translation rows are assembled into structured Elasticsearch documents.

The current synchronization pipeline has several issues:

- N+1 queries against the translation table
- No resilience to partial failures
- No visibility into stale translation data

The objective of this exercise is to design and implement a cleaner and more scalable translation loading layer.

---

# Data Model (Relevant Excerpt)

## PRODUCT

| Column      | Type    | Description        |
| ----------- | ------- | ------------------ |
| id          | uuid PK | Product UUID       |
| sku         | varchar | Product SKU        |
| part_number | varchar | Part number        |
| brand       | varchar | Brand code         |
| category_id | uuid FK | Category reference |

---

## ATTRIBUTE

| Column      | Type    | Description          |
| ----------- | ------- | -------------------- |
| id          | uuid PK | Attribute UUID       |
| code        | varchar | Attribute code       |
| metric_unit | varchar | Optional metric unit |

---

## PRODUCT_SPECIFICATION

| Column       | Type    | Description              |
| ------------ | ------- | ------------------------ |
| id           | uuid PK | Specification UUID       |
| product_id   | uuid FK | Product reference        |
| attribute_id | uuid FK | Attribute reference      |
| value        | varchar | Raw value or option code |

---

## TRANSLATION

| Column      | Type        | Description                                     |
| ----------- | ----------- | ----------------------------------------------- |
| id          | uuid PK     | Translation UUID                                |
| entity_type | varchar     | `product`, `attribute`, `product_specification` |
| entity_id   | varchar     | Stringified UUID                                |
| locale      | varchar     | `en`, `th`, `ko`                                |
| field_name  | varchar     | Translation field name                          |
| field_value | text        | Translation value                               |
| updated_at  | timestamptz | Last updated timestamp                          |

---

# What You Need To Build

## Translation Loader

Design and implement a `TranslationLoader` in Go.

The loader will be used by a product synchronization pipeline to assemble Elasticsearch documents.

### Requirements

1. Define a `TranslationLoader` interface and implement it against PostgreSQL.

2. Bulk-load translations for batches of entity IDs using a single query (no N+1 queries).

3. Support locale filtering so callers can request only specific locales.

4. Support all entity types:

   - product translations
   - attribute labels
   - specification value labels

5. Add an in-process cache layer:

   - configurable TTL
   - entity-level invalidation
   - thread-safe implementation

6. Build a `ProductDocumentBuilder` that:

   - loads product data
   - loads translations
   - assembles Elasticsearch document shape
   - gracefully handles missing translations

Missing translations must:

- never panic
- never fail document generation
- fallback to English or empty string

---

# Target Elasticsearch Document Shape

```json
{
  "uuid": "...",
  "sku": "BP-OIL-5W30-1L",
  "part_number": "5W30-1L",
  "brand": {
    "code": "bosch",
    "label": {
      "en": "Bosch",
      "th": "บอช"
    }
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
  "oil_grade": {
    "code": "5w30",
    "label": {
      "en": "5W-30",
      "th": "5W-30"
    }
  },
  "attributes": {
    "oil_grade": "5w30",
    "engine_type": "petrol"
  }
}
```

---

# Design Question (No Code Required)

How would you extend this loader to support delta sync?

Specifically:

- only reload translations updated since a given cursor timestamp
- describe the SQL/query strategy
- explain the operational trade-offs

Expected discussion areas:

- `updated_at` cursor strategy
- polling vs CDC
- eventual consistency
- stale reads
- batching trade-offs

---

# Deliverables

- Go module with clear package structure
- Translation loader implementation
- PostgreSQL-backed loader implementation
- Thread-safe TTL cache
- Docker Compose setup for PostgreSQL
- Unit tests using mocks/interfaces
- One integration test using seeded PostgreSQL data
- README with:

  - setup instructions
  - design decisions
  - trade-offs
  - future improvements
  - design question answer

---

# Evaluation Rubric

| Criterion        | What Is Evaluated                                      |
| ---------------- | ------------------------------------------------------ |
| Query Efficiency | Bulk translation loading, no N+1 queries               |
| Interface Design | Testability and abstraction quality                    |
| Cache Design     | TTL, invalidation, concurrency safety                  |
| Error Handling   | Graceful fallback and contextual errors                |
| Test Quality     | Mock-based unit tests and repeatable integration tests |
| Code Readability | Package organization and maintainability               |
| Design Reasoning | Trade-off awareness and schema understanding           |

---

# Key Engineering Focus Areas

This exercise primarily evaluates:

- scalable backend architecture
- operational thinking
- clean abstractions
- caching strategy
- testing strategy
- maintainability
- query efficiency

The emphasis is on engineering reasoning and system design quality rather than framework usage or line count.
