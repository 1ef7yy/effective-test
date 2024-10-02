package v1

import (
	"net/http"
)

func (v *Router) Groups() http.Handler {
	apimux := http.NewServeMux()

	apimux.Handle("GET /groups", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}))

	return http.StripPrefix("/api", apimux)
}
