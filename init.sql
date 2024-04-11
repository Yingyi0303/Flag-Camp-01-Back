DROP TABLE IF EXISTS replies;
DROP TABLE IF EXISTS discussions;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;

CREATE TYPE role AS ENUM ('manager', 'resident', 'third_party');
CREATE TABLE users (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role role NOT NULL
);
CREATE TABLE discussions (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    topic TEXT NOT NULL,
    content TEXT NOT NULL,
    last_update_time TIMESTAMP NOT NULL
);
CREATE TABLE replies (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    discussion_id INTEGER REFERENCES discussions(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    reply_time TIMESTAMP NOT NULL
);