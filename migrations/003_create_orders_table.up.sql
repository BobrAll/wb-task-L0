CREATE TABLE orders
(
    order_uid           TEXT PRIMARY KEY,
    track_number        TEXT        NOT NULL UNIQUE,
    entry               TEXT        NOT NULL,
    delivery_id         INT         NOT NULL REFERENCES deliveries (id),
    payment_transaction TEXT        NOT NULL REFERENCES payments (transaction),
    locale              TEXT        NOT NULL,
    internal_signature  TEXT        NOT NULL,
    customer_id         TEXT        NOT NULL,
    delivery_service    TEXT        NOT NULL,
    shard_key           VARCHAR(25) NOT NULL,
    sm_id               INT         NOT NULL,
    date_created        TIMESTAMP   NOT NULL,
    oof_shard           TEXT        NOT NULL
);