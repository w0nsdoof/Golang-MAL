package models

import (
	"database/sql"
)

type Anime struct {
	Id     int     `json:"id"`
	Rating float64 `json:"rating"`
	Title  string  `json:"title"`
	Genres string  `json:"genres"`
}

type AnimeModel struct {
	DB *sql.DB
}

func (am *AnimeModel) Insert(anime *Anime) error {
	query := `INSERT INTO animes(title, rating, genres)
			  VALUES($1, $2, $3) RETURNING id`

	err := am.DB.QueryRow(query, anime.Title, anime.Rating, anime.Genres).Scan(&anime.Id)
	if err != nil {
		return err
	}
	return nil
}

func (am *AnimeModel) Get(id int) (*Anime, error) {
	query := `SELECT id, title, rating, genres
	          FROM animes WHERE id = $1`

	var anime Anime

	err := am.DB.QueryRow(query, id).Scan(&anime.Id, &anime.Title, &anime.Rating, &anime.Genres)
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

func (am *AnimeModel) Update(anime *Anime) error {
	query := `UPDATE animes SET title=$1, rating=$2, genres=$3 WHERE id = &4`

	_, err := am.DB.Exec(query, anime.Title, anime.Rating, anime.Genres, anime.Id)
	if err != nil {
		return err
	}
	return nil
}

func (am *AnimeModel) Delete(id int) error {
	query := `DELETE FROM animes WHERE id=$1`

	_, err := am.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
