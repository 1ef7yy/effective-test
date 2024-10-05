package routes

import (
	v1 "emobile/internal/routes/v1"
	"emobile/internal/view"
	"net/http"
)

func InitRouter(view view.View) *http.ServeMux {
	mux := http.NewServeMux()
	v1 := v1.NewRouter(view)

	mux.Handle("/api/", v1.Api())

	mux.Handle("/api/songs/", v1.Songs())

	mux.Handle("/api/groups/", v1.Groups())

	return mux
}
