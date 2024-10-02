package models

import "time"

type InfoDTOReq struct {
	GroupID string `json:"group_id"`
	SongID  string `json:"song_id"`
}

type InfoDTOResponse struct {
	SongID      string    `json:"song_id"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type SongPost struct {
	GroupID     string    `json:"group_id"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type SongDTO struct {
	SongID      string    `json:"song_id"`
	GroupID     string    `json:"group_id"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}
