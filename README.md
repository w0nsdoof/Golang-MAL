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

Table user {
  id serial [primary key]
  username text
  email text
  password text
  created_at timestamp
}

Table anime {
  id serial [primary key]
  title text
  genres text
  rating float
}

// many-to-many
Table user_and_anime {
  id serial [primary key]
  user_id int
  anime_id int
  status text
  user_rating float
}

Ref: user_and_anime.user_id < users.id
Ref: user_and_anime.anime_id < animes.id
```

## Quick start

Requirements: golang, gobuffalo-soda, postgresql

- soda g config
- update information in database.yml
- soda migrate up
- go run .\cmd\myAnimeList\...
