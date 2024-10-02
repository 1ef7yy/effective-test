package domain

import (
	"context"
	"emobile/internal/models"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"fmt"
	"os"
)

type Domain interface {
	Info(group, song string) (models.InfoDTOResponse, error)
	GetGroup(group string) (models.Group, error)
	GetGroupSongs(group string, limit, offset int) ([]models.Song, error)
}

type domain struct {
	pg storage.Postgres
}

func NewDomain(log logger.Logger) Domain {
	return &domain{
		pg: *storage.NewPostgres(context.Background(), os.Getenv("POSTGRES_CONN"), log),
	}
}

func (d *domain) Info(group, song string) (models.InfoDTOResponse, error) {
	info, err := d.pg.Info(context.Background(), group, song)

	if err != nil {
		return models.InfoDTOResponse{}, err
	}

	if info.ReleaseDate.IsZero() {
		return models.InfoDTOResponse{}, fmt.Errorf("song %s from group %s not found", song, group)
	}

	return info, nil
}

func (d *domain) GetGroup(group string) (models.Group, error) {

	return d.pg.GetGroup(context.Background(), group)
}

func (d *domain) GetGroupSongs(group string, limit, offset int) ([]models.Song, error) {

	return d.pg.GetGroupSongs(context.Background(), group, limit, offset)
}
