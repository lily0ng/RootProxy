package config

type Settings struct {
	Theme          string
	DefaultProfile string
}

func DefaultSettings() *Settings {
	return &Settings{
		Theme:          "htb-dark",
		DefaultProfile: "htb-pentest",
	}
}
