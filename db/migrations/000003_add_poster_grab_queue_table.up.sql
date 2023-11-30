BEGIN;

CREATE TYPE poster_grab_status AS ENUM (
    'FINISHED',
    'WAITING',
    'IN PROGRESS',
    'FAILED'
);

CREATE TABLE IF NOT EXISTS poster_queue(
    id serial PRIMARY KEY,
    content_id integer UNIQUE NOT NULL,
    poster_id integer UNIQUE,
    status poster_grab_status not null,
    CONSTRAINT fk_content_id
        FOREIGN KEY(content_id)
            REFERENCES content(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_poster_id
        FOREIGN KEY(poster_id)
            REFERENCES poster(id)
            ON DELETE CASCADE
);

COMMIT;
