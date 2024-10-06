package view

import (
	"emobile/pkg/logger"
	"net/http"

	"emobile/internal/domain"
)

type View interface {
	GetSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	NewSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetAllSongs(log logger.Logger, w http.ResponseWriter, r *http.Request)
	EditSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	DeleteSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	NewGroup(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetGroupSongs(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetAllGroups(log logger.Logger, w http.ResponseWriter, r *http.Request)
	EditGroup(log logger.Logger, w http.ResponseWriter, r *http.Request)
}

type view struct {
	log    logger.Logger
	domain domain.Domain
}

func NewView(log logger.Logger) View {
	return &view{
		log:    log,
		domain: domain.NewDomain(log),
	}
}
