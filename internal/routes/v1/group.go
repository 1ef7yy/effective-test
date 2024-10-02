package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (v *Router) Groups() http.Handler {
	apimux := http.NewServeMux()
	apimux.Handle("GET /group", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		group := r.URL.Query().Get("group")
		if group == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request, missing group")
			return
		}
		songs, err := v.View.GetGroup(group)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, err := json.Marshal(songs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}))
	return http.StripPrefix("/api", apimux)
}
