package v1

import (
	"emobile/pkg/logger"
	"net/http"
)

func (v *Router) Songs() http.Handler {
	apimux := http.NewServeMux()

	log := logger.NewLogger(nil)

	apimux.Handle("GET /all", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.GetAllSongs(log, w, r)
	}))

	apimux.Handle("GET /get_song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.GetSong(log, w, r)
	}))

	apimux.Handle("POST /new_song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.NewSong(log, w, r)
	}))

	apimux.Handle("POST /edit_song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.EditSong(log, w, r)
	}))

	apimux.Handle("DELETE /delete_song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.DeleteSong(log, w, r)
	}))

	return http.StripPrefix("/api/songs", apimux)
}
