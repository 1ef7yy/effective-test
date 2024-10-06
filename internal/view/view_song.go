package view

import (
	"emobile/internal/errors"
	"emobile/internal/models"
	"emobile/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (v *view) GetSong(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	song := r.URL.Query().Get("song")
	group := r.URL.Query().Get("group")
	verse_offset, err := strconv.Atoi(r.URL.Query().Get("verse_offset"))

	if err != nil {
		log.Error(err.Error())
		httpResponse := errors.NewHTTPError(400, "Bad request, missing verse offset or write a number")
		w.WriteHeader(httpResponse.Status())
		fmt.Fprintf(w, httpResponse.Message())
		return
	}

	verse_limit, err := strconv.Atoi(r.URL.Query().Get("verse_limit"))

	if err != nil {
		log.Error(err.Error())
		httpResponse := errors.NewHTTPError(400, "Bad request, missing verse limit or write a number")
		w.WriteHeader(httpResponse.Status())
		fmt.Fprintf(w, httpResponse.Message())
		return
	}

	if group == "" || song == "" {
		log.Error(err.Error())
		httpResponse := errors.NewHTTPError(400, "Bad request, missing group or song")
		w.WriteHeader(httpResponse.Status())
		fmt.Fprintf(w, httpResponse.Message())
		return
	}

	Song, apierr := v.domain.GetSong(group, song, verse_offset, verse_limit)

	if Song.SongID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Printf("Song: %v\n", Song)

	if apierr != nil {
		log.Error(apierr.Error())
		w.WriteHeader(apierr.Status())
		w.Write([]byte(apierr.Message()))
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

	songs, apierr := v.domain.GetAllSongs()
	if apierr != nil {
		log.Error(apierr.Error())
		w.WriteHeader(apierr.Status())
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

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.GroupName == "" || data.SongName == "" || data.SongText == "" || data.Link == "" || data.ReleaseDate == "" {
		log.Error("Bad request, missing required field(s)")
		httpResponse := errors.NewHTTPError(400, "Bad request, missing required field(s)")
		w.WriteHeader(httpResponse.Status())
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

	songID, apierr := v.domain.NewSong(formatted_data)

	if apierr != nil {
		log.Error(apierr.Error())
		w.WriteHeader(apierr.Status())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Location", songID)
	w.WriteHeader(http.StatusCreated)

}

func (v *view) EditSong(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	var data models.EditSongReq

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.GroupName == "" || data.SongName == "" || data.SongText == "" || data.Link == "" || data.ReleaseDate.IsZero() {
		httpResponse := errors.NewHTTPError(400, "Bad request, missing required field(s)")
		w.WriteHeader(httpResponse.Status())
		fmt.Fprintf(w, httpResponse.Message())
		return
	}
}

func (v *view) DeleteSong(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	var data models.DeleteSongReq

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.GroupName == "" || data.SongName == "" {
		httpResponse := errors.NewHTTPError(400, "Bad request, missing required field(s)")
		w.WriteHeader(httpResponse.Status())
		fmt.Fprintf(w, httpResponse.Message())
		return
	}

	songID, apierr := v.domain.DeleteSong(data)

	if apierr != nil {
		log.Error(apierr.Error())
		w.WriteHeader(apierr.Status())
		return
	}

	if songID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(songID))

}
