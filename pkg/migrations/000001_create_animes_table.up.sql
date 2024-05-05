CREATE TABLE IF NOT EXISTS animes (
  id          SERIAL PRIMARY KEY,
  rating     FLOAT,
  title       TEXT NOT NULL,
  genres      TEXT
);