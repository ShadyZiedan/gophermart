-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdrawals
(
    id           serial primary key,
    user_id      integer        not null references users (id),
    number       bigint         not null,
    sum          decimal(10, 2) not null,
    processed_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE withdrawals;
-- +goose StatementEnd
