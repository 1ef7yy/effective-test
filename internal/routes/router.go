package routes

import (
	v1 "bean/internal/routes/v1"
	"bean/internal/view"
	"net/http"
)

func InitRouter(view view.View) *http.ServeMux {
	mux := http.NewServeMux()
	v1 := v1.NewRouter(view)

	mux.Handle("/api/", v1.Api())

	return mux
}
