-- +goose Up
CREATE TABLE pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    city TEXT NOT NULL CHECK (city IN ('москва', 'санкт-Петербург', 'казань'))
);

-- +goose Down
DROP TABLE IF EXISTS pvz;
