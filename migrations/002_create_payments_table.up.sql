CREATE TABLE payments
(
    transaction   TEXT PRIMARY KEY,
    request_id    TEXT    NOT NULL,
    currency      TEXT    NOT NULL,
    provider      TEXT    NOT NULL,
    amount        INT     NOT NULL,
    payment_dt    BIGINT  NOT NULL,
    bank          TEXT    NOT NULL,
    delivery_cost NUMERIC NOT NULL,
    goods_total   INT     NOT NULL,
    custom_fee    NUMERIC NOT NULL
);