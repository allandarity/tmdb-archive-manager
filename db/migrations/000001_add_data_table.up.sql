BEGIN;

CREATE TYPE content_type AS ENUM (
    'tv',
    'movie'
);

CREATE TABLE IF NOT EXISTS content(
    id serial PRIMARY KEY,
    title VARCHAR (250) NOT NULL,
    tmdb_id INTEGER UNIQUE NOT NULL,
    tmdb_popularity VARCHAR (7) NOT NULL,
    imdb_id VARCHAR(50) UNIQUE,
    imdb_popularity VARCHAR (7),
    content_type content_type not null
);

COMMIT;
