package models

import "time"

type SongDetail struct {
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type Song struct {
	SongID     string     `json:"song_id"`
	GroupID    string     `json:"group_id"`
	SongDetail SongDetail `json:"song_detail"`
}

type Group struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}
