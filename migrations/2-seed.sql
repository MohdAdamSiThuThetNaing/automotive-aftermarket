
-- ATTRIBUTE
INSERT INTO attribute (
    id,
    code,
    metric_unit
)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    'oil_grade',
    NULL
);

-- PRODUCT SPECIFICATION
INSERT INTO product_specification (
    id,
    product_id,
    attribute_id,
    value
)
VALUES (
    '33333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '5w30'
);

-- ATTRIBUTE TRANSLATIONS
INSERT INTO translation (
    id,
    entity_type,
    entity_id,
    locale,
    field_name,
    field_value,
    updated_at
)
VALUES
(
    uuid_generate_v4(),
    'attribute',
    '22222222-2222-2222-2222-222222222222',
    'en',
    'label',
    '5W-30',
    NOW()
),
(
    uuid_generate_v4(),
    'attribute',
    '22222222-2222-2222-2222-222222222222',
    'th',
    'label',
    '5W-30',
    NOW()
);

-- PRODUCT SPECIFICATION VALUE LABELS
INSERT INTO translation (
    id,
    entity_type,
    entity_id,
    locale,
    field_name,
    field_value,
    updated_at
)
VALUES
(
    uuid_generate_v4(),
    'product_specification',
    '33333333-3333-3333-3333-333333333333',
    'en',
    'value_label',
    '5W-30',
    NOW()
),
(
    uuid_generate_v4(),
    'product_specification',
    '33333333-3333-3333-3333-333333333333',
    'th',
    'value_label',
    '5W-30',
    NOW()
);
