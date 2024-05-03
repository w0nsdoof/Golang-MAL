-- add citext extension
CREATE EXTENSION IF NOT EXISTS citext;

-- Admin@kbtu.kz = admin@kbtu.kz
CREATE TABLE IF NOT EXISTS users
(
	id            BIGSERIAL PRIMARY KEY,
	created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	name          TEXT                        NOT NULL,
	email         CITEXT UNIQUE               NOT NULL,
	password_hash BYTEA                       NOT NULL,
	activated     BOOL                        NOT NULL,
	version       INTEGER                     NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS animes (
  id          SERIAL PRIMARY KEY,
  rating     FLOAT,
  title       TEXT NOT NULL,
  genres      TEXT
);

CREATE TABLE IF NOT EXISTS user_and_anime (
  id          SERIAL PRIMARY KEY,
  user_id     INT NOT NULL REFERENCES users(id),
  anime_id    INT NOT NULL REFERENCES animes(id),
  status      TEXT,
  user_rating FLOAT,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT user_anime_unique UNIQUE (user_id, anime_id)
);
