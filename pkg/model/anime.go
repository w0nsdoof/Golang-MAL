package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"final-project/pkg/validator"
)

type Anime struct {
	Id     int     `json:"id"`
	Rating float64 `json:"rating"`
	Title  string  `json:"title"`
	Genres string  `json:"genres"`
}

type AnimeModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (am *AnimeModel) GetAll(title string, filters Filters) ([]*Anime, Metadata, error) {
	query := fmt.Sprintf(
		`
			SELECT count(*) OVER(), id, rating, title, genres
			FROM animes
			WHERE (LOWER(title) = LOWER($1) OR ($1 = ''))
			ORDER BY %s %s, id ASC
			LIMIT $2 OFFSET $3
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.limit(), filters.offset()}

	rows, err := am.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			am.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var animes []*Anime
	for rows.Next() {
		var anime Anime
		err := rows.Scan(&totalRecords, &anime.Id, &anime.Rating, &anime.Title, &anime.Genres)
		if err != nil {
			return nil, Metadata{}, err
		}

		animes = append(animes, &anime)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return animes, metadata, nil
}

func (am *AnimeModel) Insert(anime *Anime) error {
	query := `INSERT INTO animes(title, rating, genres)
			  VALUES($1, $2, $3) RETURNING id`

	args := []interface{}{anime.Title, anime.Rating, anime.Genres}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan(&anime.Id)
}

func (am *AnimeModel) Get(id int) (*Anime, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, title, rating, genres
	          FROM animes WHERE id = $1`

	var anime Anime
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := am.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&anime.Id, &anime.Rating, &anime.Title, &anime.Genres)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive anime with id: %v, %w", id, err)
	}
	return &anime, nil
}

func (am *AnimeModel) Update(anime *Anime) error {
	query := `UPDATE animes SET title=$1, rating=$2, genres=$3 WHERE id = &4`

	args := []interface{}{anime.Title, anime.Rating, anime.Genres, anime.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return am.DB.QueryRowContext(ctx, query, args...).Scan()
}

func (am *AnimeModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM animes WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := am.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateAnime(v *validator.Validator, anime *Anime) {
	v.Check(anime.Title != "", "title", "must be provided")
	v.Check(len(anime.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(anime.Rating <= 10, "rating", "rating should not exceed 10")
}
