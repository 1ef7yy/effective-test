package v1

import (
	"emobile/internal/view"
	"fmt"

	"net/http"
)

type Router struct {
	View view.View
}

func NewRouter(view view.View) *Router {
	return &Router{
		View: view,
	}
}

func (v *Router) Api() http.Handler {
	apimux := http.NewServeMux()

	apimux.Handle("GET /ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	apimux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", "ok")
	}))

	return http.StripPrefix("/api", apimux)

}
