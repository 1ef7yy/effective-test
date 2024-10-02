package v1

import (
	"net/http"
)

func (v *Router) Info() http.Handler {
	apimux := http.NewServeMux()

	apimux.Handle("GET /info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}))
	return http.StripPrefix("/api", apimux)
}
