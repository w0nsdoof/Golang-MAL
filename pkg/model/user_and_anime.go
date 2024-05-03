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
	Id        int     `json:"id"`
	UserID    int     `json:"userID"`  // foreign key from users table
	AnimeID   int     `json:"animeID"` // foreign key from anime table
	Rating    float64 `json:"rating"`  // 1 single user 1 review on 1 anime
	Status    string  `json:"status"`  // "watching", "completed", "on hold"
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type UAModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *UAModel) Insert(ua *User_and_Anime) error {
	query := `INSERT INTO user_and_anime (userID, animeID, rating, status)
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	args := []interface{}{ua.UserID, ua.AnimeID, ua.Rating, ua.Status}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&ua.Id, &ua.CreatedAt, &ua.UpdatedAt)
}

func (m *UAModel) Get(id int) (*User_and_Anime, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, updated_at, userID, animeID, rating, status
		FROM user_and_anime WHERE id = $1`

	var ua User_and_Anime
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&ua.Id, &ua.CreatedAt, &ua.UpdatedAt, &ua.UserID, &ua.AnimeID, &ua.Rating, &ua.Status)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive user_and_anime with id: %v, %w", id, err)
	}
	return &ua, nil
}

func (m *UAModel) Update(ua User_and_Anime) error {
	query := `
		UPDATE user_and_anime
		SET rating = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND updated_at = $5
		RETURNING updated_at
		`
	args := []interface{}{ua.Rating, ua.Status, ua.Id, ua.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&ua.UpdatedAt)
}

func (m UAModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM user_and_anime
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateUA(v *validator.Validator, ua *User_and_Anime) {
	v.Check(ua.UserID < 1, "userID", "must be provided")
	v.Check(ua.AnimeID < 1, "animeID", "must be provided")
	v.Check(ua.Status != "", "status", "must be provided")
	v.Check(ua.Rating >= 0, "rating", "rating must be more than 0")
	v.Check(ua.Rating <= 10, "rating", "rating must be not more than 10")
}
