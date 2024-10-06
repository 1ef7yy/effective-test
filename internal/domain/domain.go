package domain

import (
	"context"
	"emobile/internal/errors"
	"emobile/internal/models"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"os"
)

type Domain interface {
	GetSong(group, song string, verse_offset, verse_limit int) (models.SongDTO, errors.APIError)
	GetAllSongs() ([]models.Song, errors.APIError)
	NewSong(data models.NewSongFormattedReq) (string, errors.APIError)
	EditSong(data models.EditSongReq) errors.APIError
	DeleteSong(data models.DeleteSongReq) (string, errors.APIError)
	NewGroup(data models.NewGroupReq) (string, errors.APIError)
	GetAllGroups() ([]models.Group, errors.APIError)
	GetGroupSongs(group string) ([]models.Song, errors.APIError)
	EditGroup(group_name models.Group) errors.APIError
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
