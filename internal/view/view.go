package view

import (
	"emobile/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"emobile/internal/domain"
	"emobile/internal/models"
)

type View interface {
	GetSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	NewSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	NewGroup(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetGroupSongs(logger.Logger, http.ResponseWriter, *http.Request)
	GetGroups(logger.Logger, http.ResponseWriter, *http.Request)
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

func (v *view) NewSong(log logger.Logger, w http.ResponseWriter, r *http.Request) {
	var data models.NewSongReq

	fmt.Printf("Data view: %v\n", data)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.GroupName == "" || data.SongName == "" || data.SongText == "" || data.Link == "" || data.ReleaseDate == "" {
		log.Error("Bad request, missing required field(s)")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request, missing required field(s)")
		return
	}

	formatted_date, err := time.Parse("2006-01-02", data.ReleaseDate)

	if err != nil {
		log.Error(fmt.Sprintf("Error parsing release date: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	formatted_data := models.NewSongFormattedReq{
		GroupName:   data.GroupName,
		SongName:    data.SongName,
		ReleaseDate: formatted_date,
		SongText:    data.SongText,
		Link:        data.Link,
	}

	songID, err := v.domain.NewSong(formatted_data)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Location", fmt.Sprintf("%s", songID))
	w.WriteHeader(http.StatusOK)

}

func (v *view) NewGroup(log logger.Logger, w http.ResponseWriter, r *http.Request) {
	var data models.NewGroupReq

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.GroupName == "" {
		log.Error("Bad request, missing required field(s)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	group_id, err := v.domain.NewGroup(data)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Location", group_id)
	w.WriteHeader(http.StatusOK)
}

func (v *view) GetGroupSongs(log logger.Logger, w http.ResponseWriter, r *http.Request) {
}

func (v *view) GetGroups(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	groups, err := v.domain.GetGroups()
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(groups)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
