package suite

import (
	"log/slog"
	"os"
	"testing"
	"url-shortener-pronetheus-consumer/internal/config"
)

type Suite struct {
	*testing.T
	Cfg *config.Config
	Log *slog.Logger
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(t *testing.T) *Suite {
	t.Parallel()

	log := setupLogger(envLocal)

	cfg := config.MustLoad()

	return &Suite{
		T: t,
		Cfg: cfg,
		Log: log,
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}