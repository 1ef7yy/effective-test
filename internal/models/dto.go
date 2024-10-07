package models

import (
	"encoding/json"
	"time"
)

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
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date,time_format:'2006-01-02',validate:'2006-01-02',omitempty"`
	SongText    string    `json:"song_text,omitempty"`
	Link        string    `json:"link,omitempty"`
}

type EditSongReq struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date" omitempty time_format:"2006-01-02" validate:"2006-01-02"`
	SongText    string    `json:"song_text,omitempty"`
	Link        string    `json:"link,omitempty"`
}

type DeleteSongReq struct {
	SongName  string `json:"song_name"`
	GroupName string `json:"group_name"`
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

func (data NewSongReq) IsValid() bool {
	return !(data.GroupName == "" || data.SongName == "")
}

func (data EditSongReq) IsValid() bool {
	return !(data.GroupName == "" || data.SongName == "")
}

func (e *EditSongReq) UnmarshalJSON(data []byte) error {
	var aux struct {
		GroupName   string `json:"group_name"`
		SongName    string `json:"song_name"`
		ReleaseDate string `json:"release_date"`
		SongText    string `json:"song_text"`
		Link        string `json:"link"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.GroupName = aux.GroupName
	e.SongName = aux.SongName
	e.SongText = aux.SongText
	e.Link = aux.Link

	// Custom unmarshalling for ReleaseDate
	releaseDate, err := time.Parse("2006-01-02", aux.ReleaseDate)
	if err != nil {
		return err
	}

	e.ReleaseDate = releaseDate

	return nil
}

func (data DeleteSongReq) IsValid() bool {
	return !(data.GroupName == "" || data.SongName == "")
}

func (data NewGroupReq) IsValid() bool {
	return !(data.GroupName == "")
}
