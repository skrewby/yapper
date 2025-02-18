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

-- ############################################################################
-- # Threads
-- ############################################################################
CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author INT NOT NULL,
    create_date  timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT fk_thread_author FOREIGN KEY (author) REFERENCES users (id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE users;
DROP TABLE threads;
