package domain

import (
	"context"
	"emobile/internal/models"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"os"
	"strings"
)

type Domain interface {
	GetSong(group, song string, verse_offset, verse_limit int) (models.SongDTO, error)
}

type domain struct {
	pg storage.Postgres
}

func NewDomain(log logger.Logger) Domain {
	return &domain{
		pg: *storage.NewPostgres(context.Background(), os.Getenv("POSTGRES_CONN"), log),
	}
}

func (d *domain) GetSong(group, song string, verse_offset, verse_limit int) (models.SongDTO, error) {

	Song, err := d.pg.GetSong(group, song)

	if err != nil {
		return models.SongDTO{}, err
	}

	var verses []string

	verses = strings.Split(Song.SongText, "\n\n")

	Song.SongText = strings.Join(verses[verse_offset:verse_offset+verse_limit], "\n\n")

	return models.SongDTO{
		SongID:      Song.SongID,
		GroupID:     Song.GroupID,
		SongName:    Song.SongName,
		ReleaseDate: Song.ReleaseDate,
		SongText:    Song.SongText,
		Link:        Song.Link,
	}, nil

}
