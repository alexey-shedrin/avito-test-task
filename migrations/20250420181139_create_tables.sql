-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email varchar UNIQUE NOT NULL,
    password varchar NOT NULL,
    user_role varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    city varchar NOT NULL,
    registration_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY,
    pvz_id UUID NOT NULL,
    status varchar NOT NULL,
    reception_datetime TIMESTAMP NOT NULL,
    FOREIGN KEY (pvz_id) REFERENCES pvz(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    reception_id UUID NOT NULL,
    product_type varchar NOT NULL,
    acceptance_datetime TIMESTAMP NOT NULL,
    FOREIGN KEY (reception_id) REFERENCES reception(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS reception CASCADE;
DROP TABLE IF EXISTS pvz CASCADE;
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
