CREATE TABLE deliveries
(
    id      SERIAL PRIMARY KEY,
    name    TEXT        NOT NULL,
    phone   TEXT        NOT NULL,
    zip     VARCHAR(15) NOT NULL,
    city    TEXT        NOT NULL,
    address TEXT        NOT NULL,
    region  TEXT        NOT NULL,
    email   TEXT        NOT NULL
);
