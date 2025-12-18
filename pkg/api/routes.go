package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/lily0ng/RootProxy/internal/proxy"
	"github.com/lily0ng/RootProxy/internal/rootproxy"
)

func RegisterRoutes(r *mux.Router, app *rootproxy.App) {
	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"active_profile": app.Profiles.Active(),
			"proxy_count":    len(app.Proxies.List()),
		})
	}).Methods(http.MethodGet)

	v1.HandleFunc("/proxy/list", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(app.Proxies.List())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/proxy/add", func(w http.ResponseWriter, r *http.Request) {
		var p proxy.Proxy
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := app.Proxies.Add(p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/profile/switch", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := app.Profiles.SetActive(body.Name); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodPost)
}
