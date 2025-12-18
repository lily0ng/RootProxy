package rootproxy

import (
	"time"

	"github.com/yourname/rootproxy/internal/cert"
	"github.com/yourname/rootproxy/internal/config"
	"github.com/yourname/rootproxy/internal/proxy"
)

type App struct {
	Proxies  *proxy.Manager
	Certs    *cert.Manager
	Profiles *config.ProfileStore
	Settings *config.Settings
}

func NewApp() *App {
	settings := config.DefaultSettings()
	proxies := proxy.NewManager()
	certs := cert.NewManager()
	profiles := config.NewProfileStore(settings.DefaultProfile)

	_ = proxies.Add(proxy.Proxy{
		Name: "HTB-Lab-TOR",
		Type: proxy.TypeSOCKS5,
		Host: "127.0.0.1",
		Port: 9050,
	})
	_ = proxies.Add(proxy.Proxy{
		Name: "Burp-Suite",
		Type: proxy.TypeHTTP,
		Host: "127.0.0.1",
		Port: 8080,
	})

	_ = profiles.Upsert(config.Profile{
		Name:      "htb-pentest",
		Chain:     []string{"HTB-Lab-TOR", "Burp-Suite"},
		UpdatedAt: time.Now().UTC(),
	})
	_ = profiles.SetActive(settings.DefaultProfile)

	return &App{
		Proxies:  proxies,
		Certs:    certs,
		Profiles: profiles,
		Settings: settings,
	}
}
