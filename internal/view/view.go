package view

import (
	"emobile/pkg/logger"

	"emobile/internal/domain"
	"emobile/internal/models"
)

type View interface {
	Info(group, song string) (models.InfoDTOResponse, error)
	GetGroup(group string) (models.Group, error)
	GetGroupSongs(group string, limit, offset int) ([]models.Song, error)
}

type view struct {
	domain domain.Domain
}

func NewView(log logger.Logger) View {
	return &view{
		domain: domain.NewDomain(log),
	}
}

func (v *view) Info(group, song string) (models.InfoDTOResponse, error) {
	info, err := v.domain.Info(group, song)

	if err != nil {
		return models.InfoDTOResponse{}, err
	}

	return info, nil

}

func (v *view) GetGroup(group string) (models.Group, error) {
	result, err := v.domain.GetGroup(group)

	if err != nil {
		return models.Group{}, err
	}

	return result, nil

}

func (v *view) GetGroupSongs(group string, limit, offset int) ([]models.Song, error) {
	songs, err := v.domain.GetGroupSongs(group, limit, offset)

	if err != nil {
		return []models.Song{}, err
	}

	return songs, nil
}
