-- +goose Up
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS products;
