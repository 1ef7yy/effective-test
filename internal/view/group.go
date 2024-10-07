package view

import (
	"emobile/internal/errors"
	"emobile/internal/models"
	"emobile/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

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
		httpResponse := errors.NewHTTPError(400, "Bad request, missing group")
		log.Error(httpResponse.Error())
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
		return
	}

	songs, apierr := v.domain.GetGroupSongs(group)

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

func (v *view) GetAllGroups(log logger.Logger, w http.ResponseWriter, r *http.Request) {

	groups, apierr := v.domain.GetAllGroups()
	if apierr != nil {
		log.Error(apierr.Error())
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

func (v *view) EditGroup(log logger.Logger, w http.ResponseWriter, r *http.Request) {
	var group models.Group

	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if group.GroupName == "" {
		httpResponse := errors.NewHTTPError(400, "Bad request, missing required field(s)")
		log.Error(httpResponse.Error())
		w.WriteHeader(httpResponse.Code)
		fmt.Fprintf(w, httpResponse.Msg)
		return
	}

	if err := v.domain.EditGroup(group); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

}
