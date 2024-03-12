package models

import (
	"database/sql"
)

type User_and_Anime struct {
	Id         int     `json:"id"`
	UserID     int     `json:"userID"`     // foreign key from users table
	AnimeID    int     `json:"animeID"`    // foreign key from anime table
	Status     string  `json:"status"`     // "watching", "completed", "on hold"
	UserRating float64 `json:"userRating"` // 1 single user 1 review on 1 anime
}

type UAModel struct {
	DB *sql.DB
}

func (ua *UAModel) Insert(ua_var *User_and_Anime) error {
	query := `INSERT INTO user_and_anime(userID, animeID, status, userRating)
	          VALUES($1, $2, $3, $4) RETURNING id`

	err := ua.DB.QueryRow(query, ua_var.UserID, ua_var.AnimeID, ua_var.Status, ua_var.UserRating).Scan(&ua_var.Id)
	if err != nil {
		return err
	}
	return nil
}

func (ua *UAModel) Get(id int) (*User_and_Anime, error) {
	query := `SELECT Id, UserID, AnimeID, Status, UserRating
	          FROM user_and_anime WHERE id=$1`

	var ua_var User_and_Anime
	err := ua.DB.QueryRow(query, id).Scan(ua_var.Id, ua_var.UserID, ua_var.AnimeID, ua_var.Status, ua_var.UserRating)
	if err != nil {
		return nil, err
	}
	return &ua_var, nil
}

func (ua *UAModel) Update(ua_var User_and_Anime) error {
	query := `UPDATE user_and_anime SET Status=$1, UserRating=$2 WHERE id=$3`
	_, err := ua.DB.Exec(query, ua_var.Status, ua_var.UserRating, ua_var.Id)
	if err != nil {
		return err
	}
	return nil
}

func (ua *UAModel) Delete(id int) error {
	query := `DELETE FROM user_and_anime WHERE id=$1`
	_, err := ua.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
