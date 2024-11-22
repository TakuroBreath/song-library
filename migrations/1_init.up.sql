CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date VARCHAR(50) NOT NULL,
    text TEXT NOT NULL,
    link VARCHAR(2048) NOT NULL
);