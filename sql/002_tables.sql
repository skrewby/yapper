-- +goose Up

-- ############################################################################
-- # User Management
-- ############################################################################
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    display_name TEXT NOT NULL,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    create_date  timestamptz NOT NULL DEFAULT current_timestamp,
    last_updated timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT uc_email UNIQUE (email)
);

-- +goose Down
DROP TABLE users;
