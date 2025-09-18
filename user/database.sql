create table "users" (
    id bigserial PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) unique NOT NULL,
    password TEXT not null ,
    role VARCHAR(20) default 'user'
);
