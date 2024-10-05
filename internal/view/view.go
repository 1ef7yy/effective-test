package view

import (
	"emobile/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"emobile/internal/domain"
	"emobile/internal/errors"
	"emobile/internal/models"
)

type View interface {
	GetSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	NewSong(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetAllSongs(log logger.Logger, w http.ResponseWriter, r *http.Request)
	NewGroup(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetGroupSongs(log logger.Logger, w http.ResponseWriter, r *http.Request)
	GetAllGroups(log logger.Logger, w http.ResponseWriter, r *http.Request)
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
		httpResponse := errors.NewHTTPError(400, "Bad request, missing verse offset or write a number")
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
		return
	}

	verse_limit, err := strconv.Atoi(r.URL.Query().Get("verse_limit"))

	if err != nil {
		log.Error(err.Error())
		httpResponse := errors.NewHTTPError(400, "Bad request, missing verse limit or write a number")
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
		return
	}

	if group == "" || song == "" {
		log.Error(err.Error())
		httpResponse := errors.NewHTTPError(400, "Bad request, missing group or song")
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
		return
	}

	Song, err := v.domain.GetSong(group, song, verse_offset, verse_limit)

	fmt.Printf("Song: %v\n", Song)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(Song)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (v *view) GetAllSongs(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	songs, err := v.domain.GetAllSongs()
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(songs)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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
		httpResponse := errors.NewHTTPError(400, "Bad request, missing required field(s)")
		w.WriteHeader(httpResponse.Code)
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
	w.WriteHeader(http.StatusCreated)

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
		httpResponse := errors.NewHTTPError(400, "Bad request, missing required field(s)")
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
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
	group := r.URL.Query().Get("group_name")

	if group == "" {
		log.Error("Bad request, missing group")
		httpResponse := errors.NewHTTPError(400, "Bad request, missing group")
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
		return
	}

	songs, err := v.domain.GetGroupSongs(group)

	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(songs)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (v *view) GetAllGroups(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	groups, err := v.domain.GetAllGroups()
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
