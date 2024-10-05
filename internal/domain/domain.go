package domain

import (
	"context"
	"emobile/internal/errors"
	"emobile/internal/models"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"fmt"
	"os"
	"strings"
)

type Domain interface {
	GetSong(group, song string, verse_offset, verse_limit int) (models.SongDTO, error)
	GetAllSongs() ([]models.Song, error)
	NewSong(data models.NewSongFormattedReq) (string, error)
	NewGroup(data models.NewGroupReq) (string, error)
	GetGroups() ([]models.Group, error)
}

type domain struct {
	pg    storage.Postgres
	redis storage.Redis
}

func NewDomain(log logger.Logger) Domain {
	return &domain{
		pg: *storage.NewPostgres(context.Background(), os.Getenv("POSTGRES_CONN"), log),
	}
}

func (d *domain) GetSong(group, song string, verse_offset, verse_limit int) (models.SongDTO, error) {

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

func (d *domain) GetAllSongs() ([]models.Song, error) {

	songs, err := d.pg.GetAllSongs()

	if err != nil {
		d.pg.Log.Error(err.Error())
		return nil, errors.NewHTTPError(500, err.Error())
	}

	return songs, nil

}

func (d *domain) NewSong(data models.NewSongFormattedReq) (string, error) {

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

	return id, nil

}

func (d *domain) NewGroup(data models.NewGroupReq) (string, error) {

	groupID, err := d.pg.NewGroup(context.Background(), data)

	if err != nil {
		d.pg.Log.Error(err.Error())
		return "", errors.NewHTTPError(500, err.Error())
	}

	return groupID, nil

}

func (d *domain) GetGroups() ([]models.Group, error) {
	groups, err := d.pg.GetGroups()
	if err != nil {
		return nil, errors.NewHTTPError(500, err.Error())
	}
	return groups, nil
}
