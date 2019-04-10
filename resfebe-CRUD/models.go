package main

type Resfebe struct {
	ImageID    int    `json:"_id"`
	Word       string `json:"word"`
	ImagePath  string `json:"imagePath"`
	Difficulty int    `json:"difficulty"`
	Category   int    `json:"category"`
	Date       string `json:"date"`
	Language   string `json:"language"`
}

type Admin struct {
	AdminID   int    `json:"_id"`
	Password  string `json:"password"`
	UserName  string `json:"userName"`
	Enabled   int    `json:"enabled"`
	LastLogin string `json:"lastLogin"`
}

type User struct {
	UserID      int    `json:"_id"`
	PlayStoreID string `json:"playStoreId"`
	Score       string `json:"score"`
}
