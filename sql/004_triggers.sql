-- +goose Up
CREATE TRIGGER update_users_time
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER update_users_time ON users;
