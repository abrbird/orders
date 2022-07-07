-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS postgres_fdw;
CREATE TABLE orders_order_shard_0 (
                              id serial PRIMARY KEY,
                              status VARCHAR NOT NULL,
                              created_at TIMESTAMP NOT NULL default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders_order_shard_0;
-- +goose StatementEnd
