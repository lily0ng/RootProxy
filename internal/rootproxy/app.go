package rootproxy

import (
	"time"

	"github.com/lily0ng/RootProxy/internal/cert"
	"github.com/lily0ng/RootProxy/internal/config"
	"github.com/lily0ng/RootProxy/internal/monitor"
	"github.com/lily0ng/RootProxy/internal/proxy"
)

type App struct {
	Proxies  *proxy.Manager
	Chains   *proxy.ChainStore
	Rotator  *proxy.Rotator
	Monitor  *monitor.Store
	Certs    *cert.Manager
	Profiles *config.ProfileStore
	Routing  *config.RoutingStore
	Security *config.SecurityStore
	Settings *config.Settings
}

func NewApp() *App {
	settings := config.DefaultSettings()
	proxies := proxy.NewManager()
	chains := proxy.NewChainStore()
	rotator := proxy.NewRotator()
	mon := monitor.NewStore()
	certs := cert.NewManager()
	profiles := config.NewProfileStore(settings.DefaultProfile)
	routing := config.NewRoutingStore()
	security := config.NewSecurityStore()

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
		Chains:   chains,
		Rotator:  rotator,
		Monitor:  mon,
		Certs:    certs,
		Profiles: profiles,
		Routing:  routing,
		Security: security,
		Settings: settings,
	}
}
