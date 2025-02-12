-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_modified_column()
    RETURNS TRIGGER AS
$$
BEGIN
    IF row (NEW.*) IS DISTINCT FROM row (OLD.*) THEN
        NEW.last_updated = now();
        RETURN NEW;
    ELSE
        RETURN OLD;
    END IF;
END;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION update_modified_column();
