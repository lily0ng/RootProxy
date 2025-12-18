package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"

	"github.com/yourname/rootproxy/internal/rootproxy"
	"github.com/yourname/rootproxy/internal/tui"
	"github.com/yourname/rootproxy/pkg/api"
)

func main() {
	var (
		profile = flag.String("profile", "", "profile name")
		apiAddr = flag.String("api", "", "start REST API server on address (e.g. 127.0.0.1:8081)")
	)
	flag.Parse()

	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	app := rootproxy.NewApp()
	if *profile != "" {
		_ = app.Profiles.SetActive(*profile)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if *apiAddr != "" {
		srv := api.NewServer(*apiAddr, app)
		go func() {
			if err := srv.Start(ctx); err != nil {
				logrus.WithError(err).Error("api server stopped")
			}
		}()
	}

	p := tea.NewProgram(tui.NewModel(app), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logrus.WithError(err).Fatal("rootproxy exited with error")
	}
}
