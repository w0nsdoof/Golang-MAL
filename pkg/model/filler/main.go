package filler

import (
	model "final-project/pkg/model"
)

func PopulateDatabase(models model.Models) error {
	for _, anime := range animes {
		models.Animes.Insert(&anime)
	}
	// TODO: Implement users? pupulation
	// TODO: Implement user_and_anime pupulation
	return nil
}

var animes = []model.Anime{
	{Title: "Jujutsu Kaisen", Rating: 9.2, Genres: "Adventure, Dark Fantasy, Supernatural"},
	{Title: "Attack on Titan", Rating: 9.5, Genres: "Action, Military, Mystery"},
	{Title: "My Hero Academia", Rating: 8.5, Genres: "Action, Super Power, School"},
	{Title: "Demon Slayer", Rating: 9.3, Genres: "Action, Historical, Supernatural"},
	{Title: "One Piece", Rating: 9.1, Genres: "Action, Adventure, Fantasy"},
	{Title: "Tokyo Ghoul", Rating: 8.7, Genres: "Action, Mystery, Horror"},
	{Title: "Death Note", Rating: 9.0, Genres: "Mystery, Police, Psychological"},
	{Title: "Steins;Gate", Rating: 9.2, Genres: "Thriller, Sci-Fi"},
	{Title: "Fullmetal Alchemist: Brotherhood", Rating: 9.4, Genres: "Action, Military, Adventure"},
	{Title: "Hunter x Hunter", Rating: 9.1, Genres: "Action, Adventure, Fantasy"},
	{Title: "Naruto Shippuden", Rating: 8.6, Genres: "Action, Adventure, Martial Arts"},
	{Title: "Sword Art Online", Rating: 7.5, Genres: "Action, Game, Adventure"},
	{Title: "One Punch Man", Rating: 8.8, Genres: "Action, Sci-Fi, Comedy"},
	{Title: "Re:Zero", Rating: 8.5, Genres: "Drama, Fantasy, Psychological"},
	{Title: "Mob Psycho 100", Rating: 8.9, Genres: "Action, Slice of Life, Comedy"},
	{Title: "Your Name", Rating: 9.2, Genres: "Drama, Romance, Supernatural"},
	{Title: "Neon Genesis Evangelion", Rating: 8.9, Genres: "Action, Dementia, Mecha"},
	{Title: "Cowboy Bebop", Rating: 8.9, Genres: "Action, Adventure, Space"},
	{Title: "Fate/Stay Night: Unlimited Blade Works", Rating: 8.5, Genres: "Action, Supernatural, Magic"},
	{Title: "The Promised Neverland", Rating: 8.7, Genres: "Sci-Fi, Mystery, Horror"},
	{Title: "Bleach", Rating: 7.9, Genres: "Action, Adventure, Supernatural"},
}
