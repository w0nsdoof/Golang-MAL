# MyAnimeList

by Zhumabayev Askar 22B030361

## Introduction

MyAnimeList is a application that provides its users with the ability to organize, save and rate anime.

## Users REST API

```
POST /api/v1/users
GET /api/v1/users/{id}
PUT /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

## DB Structure

## ![db_schema](assets/images/db_schema.png)

```
// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Table users {
  id            bigint [pk]
  created_at    timestamp
  name          text
  email         citext [unique]
  password_hash bytea
  activated     bool
  version       integer
}

Table animes {
  id      serial [pk]
  rating  float
  title   text
  genres  text
}

Table user_and_anime {
  id          serial [pk]
  userID      int [not null]
  animeID     int [not null]
  rating      float
  review      text
  created_at  timestamp
  updated_at  timestamp
}

Ref: user_and_anime.userID > users.id
Ref: user_and_anime.animeID > animes.id
```

## How to run app

### Using app golang directly on Terminal

Provide all needed correct values.

```shell
soda migrate up
go run ./cmd/myAnimeList/
-dsn="postgres://postgres:password@localhost:5432/database_name?sslmode=disable" \
-migrations=migrations \
-fill=true \
-env=development \
-port=8081
```

### Run with docker-compose

Up database service

```bash
  docker compose up -d db
```

Build and up myanimelist service

```bash
  docker compose build myanimelist
  docker compose up myanimelist
```
