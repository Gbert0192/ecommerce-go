create table product_category(
    id serial primary key,
    name varchar(255) unique not null
)

create table product (
    id bigserial PRIMARY KEY,
    name varchar(255) not null,
    description text,
    price numeric not null,
    stock integer not null,
    category_id integer not null,
    constraint fk_category foreign key (category_id) REFERENCES product_category(id) ON DELETE CASCADE
);
