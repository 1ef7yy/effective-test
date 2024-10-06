package domain

import (
	"emobile/internal/errors"
	"emobile/internal/models"
	"fmt"
	"strings"
)

func (d *domain) GetSong(group, song string, verse_offset, verse_limit int) (models.SongDTO, errors.APIError) {

	// cached_song, err := d.redis.GetSong(group, song)

	// if err != nil {
	// 	return models.SongDTO{}, err
	// }

	// if cached_song.SongID != "" {

	// 	Song := models.SongDTO{
	// 		SongID:      cached_song.SongID,
	// 		GroupID:     cached_song.GroupID,
	// 		SongName:    cached_song.SongName,
	// 		ReleaseDate: cached_song.ReleaseDate,
	// 		SongText:    cached_song.SongText,
	// 		Link:        cached_song.Link,
	// 	}

	// 	return Song, nil
	// }

	Song, err := d.pg.GetSong(group, song)

	if err != nil {
		d.pg.Log.Error(err.Error())
		return models.SongDTO{}, errors.NewHTTPError(500, err.Error())
	}

	if Song.SongID == "" {
		return models.SongDTO{}, errors.NewHTTPError(404, "Song not found")
	}

	var verses []string

	verses = strings.Split(Song.SongText, "\n\n")

	if verse_limit > len(verses)-verse_offset {
		return models.SongDTO{}, errors.NewHTTPError(400, "Bad request, verse_offset or limit out of bounds")
	}

	return models.SongDTO{
		SongID:      Song.SongID,
		GroupID:     Song.GroupID,
		SongName:    Song.SongName,
		ReleaseDate: Song.ReleaseDate,
		SongText:    strings.Join(verses[verse_offset:verse_offset+verse_limit], "\n\n"),
		Link:        Song.Link,
	}, nil

}

func (d *domain) GetAllSongs() ([]models.Song, errors.APIError) {

	songs, err := d.pg.GetAllSongs()

	if err != nil {
		d.pg.Log.Error(err.Error())
		return nil, errors.NewHTTPError(500, err.Error())
	}

	if len(songs) == 0 {
		return nil, errors.NewHTTPError(404, "No songs found")
	}

	return songs, nil

}

func (d *domain) NewSong(data models.NewSongFormattedReq) (string, errors.APIError) {

	fmt.Printf("Data domain: %v\n", data)

	if data.GroupName == "" || data.SongName == "" || data.SongText == "" || data.Link == "" || data.ReleaseDate.IsZero() {
		return "", errors.NewHTTPError(400, "Bad request, missing required field(s)")
	}

	SongReq := models.NewSongReq{
		GroupName:   data.GroupName,
		SongName:    data.SongName,
		ReleaseDate: data.ReleaseDate.Format("2006-01-02"),
		SongText:    data.SongText,
		Link:        data.Link,
	}

	id, err := d.pg.NewSong(SongReq)

	if err != nil {
		d.pg.Log.Error(err.Error())
		return "", errors.NewHTTPError(500, err.Error())
	}

	if id == "" {
		return "", errors.NewHTTPError(404, "Not found")
	}

	return id, nil

}

func (d *domain) EditSong(data models.EditSongReq) errors.APIError {

	songID, err := d.pg.EditSong(data)

	if err != nil {
		d.pg.Log.Error(err.Error())
		return errors.NewHTTPError(500, err.Error())
	}

	if songID == "" {
		return errors.NewHTTPError(404, "Not found")
	}

	return nil
}

func (d *domain) DeleteSong(data models.DeleteSongReq) (string, errors.APIError) {
	songID, err := d.pg.DeleteSong(data)

	if songID != "" {
		return "", errors.NewHTTPError(404, "Not found")
	}

	if err != nil {
		d.pg.Log.Error(err.Error())
		return "", errors.NewHTTPError(500, err.Error())
	}
	return songID, nil
}
