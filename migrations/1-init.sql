CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR NOT NULL,
    part_number VARCHAR NOT NULL,
    brand VARCHAR NOT NULL,
    category_id UUID
);

CREATE TABLE attribute (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR NOT NULL UNIQUE,
    metric_unit VARCHAR
);

CREATE TABLE product_specification (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES product(id),
    attribute_id UUID NOT NULL REFERENCES attribute(id),
    value VARCHAR NOT NULL
);

CREATE TABLE translation (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR NOT NULL,
    entity_id VARCHAR NOT NULL,
    locale VARCHAR NOT NULL,
    field_name VARCHAR NOT NULL,
    field_value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_translation_lookup
ON translation(entity_type, entity_id, locale);

CREATE UNIQUE INDEX idx_translation_unique
ON translation(entity_type, entity_id, locale, field_name);