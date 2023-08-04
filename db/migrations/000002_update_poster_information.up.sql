BEGIN;

CREATE TABLE IF NOT EXISTS poster(
    id serial PRIMARY KEY,
    content_id integer UNIQUE NOT NULL,
    poster_url VARCHAR(50),
    poster_data bytea,
    CONSTRAINT fk_customer
        FOREIGN KEY(content_id)
            REFERENCES content(id)
            ON DELETE CASCADE
);

COMMIT;
