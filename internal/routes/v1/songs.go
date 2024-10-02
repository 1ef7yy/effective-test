package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (v *Router) Songs() http.Handler {
	apimux := http.NewServeMux()

	apimux.Handle("POST /song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	apimux.Handle("GET /song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		song := r.URL.Query().Get("song")
		if group == "" || song == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request, missing group or song")
			return
		}
		info, err := v.View.Info(group, song)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if info.SongName == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Song %s of group %s not found", song, group)))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, err := json.Marshal(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}))

	apimux.Handle("GET /songs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request, missing offset or offset not integer")
			return
		}
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request, missing limit or limit not integer")
			return
		}
		if group == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request, missing group")
			return
		}
		songs, err := v.View.GetGroupSongs(group, offset, limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, err := json.Marshal(songs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}))

	return http.StripPrefix("/api", apimux)
}
