-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE table if not exists "metrics"
(
    "name"  varchar UNIQUE NOT NULL,
    "type"  varchar        NOT null,
    "delta" integer,
    "value" decimal,
    "hash"  varchar
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "metrics";
