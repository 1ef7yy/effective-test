package view

import (
	"emobile/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"emobile/internal/domain"
)

type View interface {
	GetSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
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

func (v *view) GetSong(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	song := r.URL.Query().Get("song")
	group := r.URL.Query().Get("group")
	verse_offset, err := strconv.Atoi(r.URL.Query().Get("verse_offset"))

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request, missing verse offset or write a number")
		return
	}

	verse_limit, err := strconv.Atoi(r.URL.Query().Get("verse_limit"))

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request, missing verse limit or write a number")
		return
	}

	if group == "" || song == "" {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request, missing group or song")
		return
	}

	info, err := v.domain.GetSong(group, song, verse_offset, verse_limit)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(info)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Song, err := v.domain.GetSong(group, song, verse_limit, verse_offset)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err = json.Marshal(Song)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
