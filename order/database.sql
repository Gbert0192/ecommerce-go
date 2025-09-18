create table order_detail(
    id bigserial primary key,
    products text not null ,
    order_history text not null
)

create table orders (
    id bigserial primary key ,
    user_id bigint not null,
    amount numeric not null,
    total_qty integer not null,
    payment_method varchar(50),
    shipping_address text,
    status integer not null,
    order_detail_id bigint references order_detail(id),
    create_time timestamp default current_timestamp,
    update_time timestamp default current_timestamp
)

create table order_request_log(
    id bigserial primary key,
    idempotency_token text unique not null,
    crate_time timestamp default current_timestamp
)
