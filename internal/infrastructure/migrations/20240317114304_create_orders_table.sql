-- +goose Up
-- +goose StatementBegin
create type order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
create table orders
(
    id           serial primary key,
    user_id      integer      not null references users (id),
    number       bigint       not null,
    status       order_status not null,
    accrual      decimal(10, 2),
    uploaded_at  timestamp    not null default current_timestamp,
    processed_at timestamp
);
create index orders_user_id_idx on orders (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
DROP TYPE order_status;
-- +goose StatementEnd
