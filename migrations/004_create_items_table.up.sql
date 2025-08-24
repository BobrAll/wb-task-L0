CREATE TABLE items
(
    chrt_id      BIGINT PRIMARY KEY,
    track_number TEXT    NOT NULL REFERENCES orders (track_number),
    price        NUMERIC NOT NULL,
    rid          TEXT    NOT NULL,
    name         TEXT    NOT NULL,
    sale         REAL    NOT NULL,
    size         TEXT    NOT NULL,
    total_price  NUMERIC NOT NULL,
    nm_id        BIGINT  NOT NULL,
    brand        TEXT    NOT NULL,
    status       INT     NOT NULL
);