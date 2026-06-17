-- +goose Up
ALTER TABLE couriers
ADD COLUMN IF NOT EXISTS transport_type TEXT NOT NULL DEFAULT 'on_foot';  -- on_foot | scooter | car

-- +goose Down
ALTER TABLE couriers
DROP COLUMN IF EXISTS transport_type;
