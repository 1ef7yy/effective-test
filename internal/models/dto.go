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

type SongDTO struct {
	SongID      string    `json:"song_id"`
	GroupID     string    `json:"group_id"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type NewSongReq struct {
	GroupName   string `json:"group_name"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	SongText    string `json:"song_text"`
	Link        string `json:"link"`
}

type NewSongFormattedReq struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type RedisSongDTO struct {
	SongID      string    `json:"song_id"`
	GroupID     string    `json:"group_id"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}

type NewGroupReq struct {
	GroupName string `json:"group_name"`
}

type SongDB struct {
	SongID      string    `json:"song_id"`
	GroupID     string    `json:"group_id"`
	ReleaseDate time.Time `json:"release_date"`
	SongName    string    `json:"song_name"`
	SongText    string    `json:"song_text"`
	Link        string    `json:"link"`
}
