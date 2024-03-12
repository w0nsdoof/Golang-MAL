package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (username, email, password, created_at)
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := um.DB.QueryRow(query, user.Username, user.Email, user.Password, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Get(id int) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`

	/*
		SELECT u.username, u.email,u.created_At a.id, a.title, a.genres, ua.user_rating
		FROM users AS u
		INNER JOIN user_and_anime AS ua ON u.id = ua.user_id
		INNER JOIN animes AS a ON ua.anime_id = a.id
	*/

	var user User
	err := um.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserModel) Update(user *User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	_, err := um.DB.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := um.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
