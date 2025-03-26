package main

import "time"

type user struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type wine struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Producer  string    `json:"producer"`
	Country   string    `json:"country"`
	Vintage   int       `json:"vintage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type rating struct {
	ID        string    `json:"id"`
	WineID    string    `json:"wine_id"`
	UserID    string    `json:"user_id"`
	Rating    string    `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
