package v1

import (
	"emobile/pkg/logger"
	"net/http"
)

func (v *Router) Info() http.Handler {
	apimux := http.NewServeMux()

	log := logger.NewLogger(nil)

	apimux.Handle("GET /info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.GetSong(log, w, r)
	}))
	return http.StripPrefix("/api", apimux)
}
