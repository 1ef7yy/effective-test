package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (v *Router) Info() http.Handler {
	apimux := http.NewServeMux()

	apimux.Handle("GET /info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, err := json.Marshal(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}))
	return http.StripPrefix("/api", apimux)
}
