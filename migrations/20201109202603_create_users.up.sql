CREATE TABLE users (
    id serial not null primary key,
    email text not null unique,
    encrypted_password text not null
);