package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/lily0ng/RootProxy/internal/cert"
	"github.com/lily0ng/RootProxy/internal/config"
	"github.com/lily0ng/RootProxy/internal/proxy"
	"github.com/lily0ng/RootProxy/internal/rootproxy"
)

func RegisterRoutes(r *mux.Router, app *rootproxy.App) {
	v1 := r.PathPrefix("/api/v1").Subrouter()
	writeJSON := func(w http.ResponseWriter, status int, v any) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(v)
	}
	writeErr := func(w http.ResponseWriter, status int, err error) {
		http.Error(w, err.Error(), status)
	}

	v1.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		activeProxy := app.Proxies.ActiveName()
		writeJSON(w, http.StatusOK, map[string]any{
			"active_profile": app.Profiles.Active(),
			"active_proxy":   activeProxy,
			"proxy_count":    len(app.Proxies.List()),
			"chain_count":    len(app.Chains.List()),
		})
	}).Methods(http.MethodGet)

	v1.HandleFunc("/proxy/list", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Proxies.List())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/proxy/active", func(w http.ResponseWriter, _ *http.Request) {
		p, ok := app.Proxies.GetActive()
		if !ok {
			writeJSON(w, http.StatusOK, map[string]any{"active": nil})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"active": p})
	}).Methods(http.MethodGet)

	v1.HandleFunc("/proxy/active", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Proxies.SetActive(body.Name); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"active": body.Name})
	}).Methods(http.MethodPost)

	v1.HandleFunc("/proxy/add", func(w http.ResponseWriter, r *http.Request) {
		var p proxy.Proxy
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Proxies.Add(p); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, p)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/proxy/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var p proxy.Proxy
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Proxies.Update(id, p); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		updated, _ := app.Proxies.GetByID(id)
		writeJSON(w, http.StatusOK, updated)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/proxy/remove/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if err := app.Proxies.Remove(id); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)

	v1.HandleFunc("/proxy/test", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = app.Proxies.ActiveName()
		}
		p, ok := app.Proxies.GetByName(name)
		if !ok {
			writeErr(w, http.StatusBadRequest, errors.New("proxy not found"))
			return
		}
		tmo := 3 * time.Second
		if qs := r.URL.Query().Get("timeout_ms"); qs != "" {
			if ms, err := strconv.Atoi(qs); err == nil && ms > 0 {
				tmo = time.Duration(ms) * time.Millisecond
			}
		}
		ctx, cancel := context.WithTimeout(r.Context(), tmo)
		defer cancel()
		tr := proxy.TestConnectivity(ctx, p)
		app.Monitor.RecordTest(name, tr)
		writeJSON(w, http.StatusOK, tr)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/proxy/export", func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}
		items := app.Proxies.List()
		switch format {
		case "json":
			b, err := proxy.ExportJSON(items)
			if err != nil {
				writeErr(w, http.StatusInternalServerError, err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)
		case "text":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(proxy.ExportText(items))
		default:
			writeErr(w, http.StatusBadRequest, errors.New("unsupported export format"))
		}
	}).Methods(http.MethodGet)

	v1.HandleFunc("/proxy/import", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		format := r.URL.Query().Get("format")
		if format == "" {
			format = proxy.DetectImportFormat(body, r.Header.Get("Content-Type"))
		}

		var parsed []proxy.Proxy
		switch format {
		case "json":
			parsed, err = proxy.ImportJSON(bytes.NewReader(body))
		case "text":
			parsed, err = proxy.ImportText(bytes.NewReader(body))
		default:
			err = errors.New("unsupported import format")
		}
		if err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}

		type failure struct {
			Name  string `json:"name"`
			Error string `json:"error"`
		}
		res := struct {
			Imported int       `json:"imported"`
			Added    int       `json:"added"`
			Failed   []failure `json:"failed"`
		}{Imported: len(parsed)}

		for _, p := range parsed {
			if p.Type == "" {
				p.Type = proxy.TypeHTTP
			}
			if p.Auth == "" {
				p.Auth = proxy.AuthNone
			}
			if err := app.Proxies.Add(p); err != nil {
				res.Failed = append(res.Failed, failure{Name: p.Name, Error: err.Error()})
				continue
			}
			res.Added++
		}
		writeJSON(w, http.StatusOK, res)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/profile/switch", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Profiles.SetActive(body.Name); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"active": body.Name})
	}).Methods(http.MethodPost)

	v1.HandleFunc("/profile/list", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Profiles.List())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/profile/upsert", func(w http.ResponseWriter, r *http.Request) {
		var p config.Profile
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Profiles.Upsert(p); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, p)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/chain/list", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Chains.List())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/chain/upsert", func(w http.ResponseWriter, r *http.Request) {
		var c proxy.Chain
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Chains.Upsert(c, 5); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, c)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/chain/remove/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		if err := app.Chains.Remove(name); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)

	v1.HandleFunc("/routing/list", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Routing.List())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/routing/upsert", func(w http.ResponseWriter, r *http.Request) {
		var rr config.RoutingRule
		if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Routing.Upsert(rr); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, rr)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/routing/remove/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if err := app.Routing.Remove(id); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)

	v1.HandleFunc("/rotation/rotate", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Profile string `json:"profile"`
			Mode    string `json:"mode"`
			Enabled bool   `json:"enabled"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Profile == "" {
			body.Profile = app.Profiles.Active()
		}
		if body.Mode == "" {
			body.Mode = string(config.RotationRoundRobin)
		}
		policy := config.RotationPolicy{Enabled: body.Enabled, Mode: config.RotationMode(body.Mode)}
		if !policy.Enabled {
			policy.Enabled = true
		}

		var chain []string
		for _, p := range app.Profiles.List() {
			if p.Name == body.Profile {
				chain = p.Chain
				break
			}
		}
		name, err := app.Rotator.Rotate(body.Profile, chain, policy, app.Proxies)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"active": name})
	}).Methods(http.MethodPost)

	v1.HandleFunc("/cert/list", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Certs.List())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/cert/add", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
			PEM  []byte `json:"pem"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		if err := app.Certs.Add(body.Name, body.PEM); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"name": body.Name})
	}).Methods(http.MethodPost)

	v1.HandleFunc("/cert/generate_self_signed", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name          string `json:"name"`
			CommonName    string `json:"common_name"`
			ValidForHours int    `json:"valid_for_hours"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Name == "" {
			body.Name = "RootProxy"
		}
		valid := 365 * 24 * time.Hour
		if body.ValidForHours > 0 {
			valid = time.Duration(body.ValidForHours) * time.Hour
		}
		gen, err := cert.GenerateSelfSigned(cert.SelfSignedOptions{CommonName: body.CommonName, ValidFor: valid})
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err)
			return
		}
		_ = app.Certs.Add(body.Name, gen.CertPEM)
		writeJSON(w, http.StatusOK, map[string]any{"name": body.Name, "cert_pem": gen.CertPEM, "key_pem": gen.KeyPEM})
	}).Methods(http.MethodPost)

	v1.HandleFunc("/security/get", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Security.Get())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/security/set", func(w http.ResponseWriter, r *http.Request) {
		var s config.SecuritySettings
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		app.Security.Set(s)
		writeJSON(w, http.StatusOK, s)
	}).Methods(http.MethodPost)

	v1.HandleFunc("/monitoring/metrics", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, app.Monitor.Snapshot())
	}).Methods(http.MethodGet)

	v1.HandleFunc("/monitoring/started", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"started_at": app.Monitor.StartedAt()})
	}).Methods(http.MethodGet)

	v1.HandleFunc("/integrations/burp/env", func(w http.ResponseWriter, _ *http.Request) {
		p, ok := app.Proxies.GetActive()
		if !ok {
			writeJSON(w, http.StatusOK, map[string]any{"HTTP_PROXY": "", "HTTPS_PROXY": ""})
			return
		}
		addr := "http://" + p.Address()
		if p.Type == proxy.TypeHTTPS {
			addr = "https://" + p.Address()
		}
		writeJSON(w, http.StatusOK, map[string]any{"HTTP_PROXY": addr, "HTTPS_PROXY": addr})
	}).Methods(http.MethodGet)

	v1.HandleFunc("/integrations/proxychains/conf", func(w http.ResponseWriter, r *http.Request) {
		profileName := r.URL.Query().Get("profile")
		if profileName == "" {
			profileName = app.Profiles.Active()
		}
		var chain []string
		for _, p := range app.Profiles.List() {
			if p.Name == profileName {
				chain = p.Chain
				break
			}
		}
		var b bytes.Buffer
		b.WriteString("strict_chain\n")
		b.WriteString("proxy_dns\n")
		b.WriteString("[ProxyList]\n")
		for _, name := range chain {
			px, ok := app.Proxies.GetByName(name)
			if !ok {
				continue
			}
			scheme := string(px.Type)
			if scheme == string(proxy.TypeHTTPS) {
				scheme = string(proxy.TypeHTTP)
			}
			b.WriteString(scheme + " " + px.Host + " " + strconv.Itoa(px.Port) + "\n")
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b.Bytes())
	}).Methods(http.MethodGet)
}
