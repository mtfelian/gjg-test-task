CREATE TABLE IF NOT EXISTS levels
(
    id   UUID PRIMARY KEY,
    x    INT   NOT NULL DEFAULT 0,
    y    INT   NOT NULL DEFAULT 0,
    maze BYTEA NOT NULL
);



