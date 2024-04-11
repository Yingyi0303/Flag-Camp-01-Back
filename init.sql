DROP TRIGGER IF EXISTS update_bill_trigger ON bills;
DROP TRIGGER IF EXISTS update_payment_trigger ON payments;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS maintenances;
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
    subject TEXT NOT NULL,
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
CREATE TABLE maintenances (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    subject TEXT NOT NULL,
    content TEXT NOT NULL,
    reply TEXT DEFAULT '',
    completed BOOLEAN DEFAULT FALSE,
    last_update_time TIMESTAMP NOT NULL
);
CREATE TABLE bills (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    maintenance_id INTEGER REFERENCES maintenances(id) ON DELETE CASCADE,
    item TEXT NOT NULL,
    amount INTEGER NOT NULL,
    bill_time TIMESTAMP NOT NULL
);
CREATE TABLE payments (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    item TEXT NOT NULL,
    amount INTEGER NOT NULL,
    payment_time TIMESTAMP NOT NULL
);
CREATE TABLE balances (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE CASCADE,
    amount INTEGER DEFAULT 0
);
CREATE OR REPLACE FUNCTION update_bill()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE balances SET amount = amount + NEW.amount WHERE username = NEW.username;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_bill_trigger AFTER INSERT ON bills FOR EACH ROW EXECUTE FUNCTION update_bill();
CREATE OR REPLACE FUNCTION update_payment()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE balances SET amount = amount - NEW.amount WHERE username = NEW.username;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_payment_trigger AFTER INSERT ON payments FOR EACH ROW EXECUTE FUNCTION update_payment();