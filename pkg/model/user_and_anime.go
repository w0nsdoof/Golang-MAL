package model

import (
	"context"
	"database/sql"
	"final-project/pkg/validator"
	"fmt"
	"log"
	"time"
)

type User_and_Anime struct {
	Id         int     `json:"id"`
	UserID     int     `json:"userID"`   // foreign key from users table
	Username   string  `json:"username"` //
	AnimeID    int     `json:"animeID"`  // foreign key from anime table
	AnimeTitle string  `json:"animeTitle"`
	Rating     float64 `json:"rating"` // 1 single user 1 review on 1 anime
	Review     string  `json:"review"` // "watching", "completed", "on hold"
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type UAModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *UAModel) AverageRating(animeID int) (float64, int, string, error) {
	var avgRating float64
	var userCount int
	var title string

	query := `
        SELECT COALESCE(AVG(ua.rating), 0), COUNT(ua.userID), a.title
        FROM user_and_anime AS ua
        JOIN animes AS a ON ua.animeID = a.id
        WHERE ua.animeID = $1
        GROUP BY a.title
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, animeID).Scan(&avgRating, &userCount, &title)
	if err != nil {
		return 0, 0, "", err
	}

	return avgRating, userCount, title, nil
}

func (m *UAModel) GetAllByUser(userID int, filters Filters) ([]*User_and_Anime, Metadata, error) {
	query := fmt.Sprintf(
		`
			SELECT count(*) OVER(), ua.id, ua.created_at, ua.updated_at, ua.userID, u.name, ua.animeID, a.title, ua.rating, ua.review
			FROM user_and_anime AS ua
			JOIN animes AS a ON ua.animeID = a.id
			JOIN users AS u ON ua.userID = u.id
			WHERE ua.userID = $1
			ORDER BY %s %s, ua.id ASC
			LIMIT $2 OFFSET $3
		`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{userID, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0
	var userAnimes []*User_and_Anime
	for rows.Next() {
		var ua User_and_Anime
		err := rows.Scan(&totalRecords, &ua.Id, &ua.CreatedAt, &ua.UpdatedAt, &ua.UserID, &ua.Username, &ua.AnimeID, &ua.AnimeTitle, &ua.Rating, &ua.Review)
		if err != nil {
			return nil, Metadata{}, err
		}

		userAnimes = append(userAnimes, &ua)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return userAnimes, metadata, nil
}

func (m *UAModel) Insert(ua *User_and_Anime) error {
	query := `INSERT INTO user_and_anime (userID, animeID, rating, review)
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	args := []interface{}{ua.UserID, ua.AnimeID, ua.Rating, ua.Review}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&ua.Id, &ua.CreatedAt, &ua.UpdatedAt)
}

func (m *UAModel) Get(id int) (*User_and_Anime, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT ua.id, ua.created_at, ua.updated_at, ua.userID, u.name, ua.animeID, a.title, ua.rating, ua.review
		FROM user_and_anime AS ua
		JOIN animes AS a ON ua.animeID = a.id
		JOIN users AS u ON ua.userID = u.id
		WHERE ua.id = $1
	`

	var ua User_and_Anime
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&ua.Id, &ua.CreatedAt, &ua.UpdatedAt, &ua.UserID, &ua.Username, &ua.AnimeID, &ua.AnimeTitle, &ua.Rating, &ua.Review)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("cannot retrieve user_and_anime with id: %v, %w", id, err)
	}

	return &ua, nil
}

func (m *UAModel) Update(ua User_and_Anime) error {
	query := `
		UPDATE user_and_anime
		SET rating = $1, review = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND updated_at = $4
		RETURNING updated_at
		`
	args := []interface{}{ua.Rating, ua.Review, ua.Id, ua.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&ua.UpdatedAt)
}

func (m UAModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM user_and_anime WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateUA(v *validator.Validator, ua *User_and_Anime) {
	v.Check(ua.UserID > 0, "userID", "must be provided")
	v.Check(ua.AnimeID > 0, "animeID", "must be provided")
	v.Check(ua.Rating >= 0, "rating", "rating must be more than 0")
	v.Check(ua.Rating <= 10, "rating", "rating must be not more than 10")
}

// UserExists checks if a user with the given ID exists.
func (m *UAModel) UserExists(userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// AnimeExists checks if an anime with the given ID exists.
func (m *UAModel) AnimeExists(animeID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM animes WHERE id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, animeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
