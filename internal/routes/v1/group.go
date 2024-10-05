package v1

import (
	"emobile/pkg/logger"
	"net/http"
)

func (v *Router) Groups() http.Handler {
	apimux := http.NewServeMux()
	log := logger.NewLogger(nil)

	// root/api/groups
	// body: { "group_name": "group_name" }
	// returning: uuid of requested group
	apimux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.GetAllGroups(log, w, r)
	}))

	apimux.Handle("GET /songs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.GetGroupSongs(log, w, r)
	}))

	apimux.Handle("POST /new_group", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.View.NewGroup(log, w, r)
	}))
	return http.StripPrefix("/api/groups", apimux)
}
