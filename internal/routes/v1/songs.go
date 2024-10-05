package v1

import (
	"emobile/pkg/logger"
	"net/http"
)

func (v *Router) Songs() http.Handler {
	apimux := http.NewServeMux()

	log := logger.NewLogger(nil)

	apimux.Handle("GET /song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.GetSong(log, w, r)
	}))

	apimux.Handle("POST /new_song", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v.View.NewSong(log, w, r)
	}))

	return http.StripPrefix("/api/songs", apimux)
}
