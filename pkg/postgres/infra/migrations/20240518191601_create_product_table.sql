-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product (
	id          UUID PRIMARY KEY,
	name        VARCHAR (127),
	description VARCHAR (511),
	created_at  TIMESTAMP,
	updated_at  TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product;
-- +goose StatementEnd
