CREATE TABLE IF NOT EXISTS users (
  id          SERIAL PRIMARY KEY,
  username    TEXT UNIQUE NOT NULL,
  email       TEXT UNIQUE NOT NULL,
  password    TEXT NOT NULL,
  created_at  TIMESTAMP NOT NULL DEFAULT now()
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
  CONSTRAINT user_anime_unique UNIQUE (user_id, anime_id) 
);