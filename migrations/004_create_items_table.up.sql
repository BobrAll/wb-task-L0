CREATE TABLE items
(
    chrt_id      BIGINT PRIMARY KEY,
    order_id     TEXT    NOT NULL REFERENCES orders (order_uid),
    track_number TEXT    NOT NULL,
    price        NUMERIC NOT NULL,
    rid          TEXT    NOT NULL,
    name         TEXT    NOT NULL,
    sale         REAL    NOT NULL,
    size         INT     NOT NULL,
    total_price  NUMERIC NOT NULL,
    nm_id        BIGINT  NOT NULL,
    brand        TEXT    NOT NULL,
    status       INT     NOT NULL
);

