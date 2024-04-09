DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;

CREATE TYPE role AS ENUM ('manager', 'resident', 'third_party');
CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    role role NOT NULL
);