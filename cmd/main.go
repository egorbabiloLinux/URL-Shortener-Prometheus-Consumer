package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener-pronetheus-consumer/internal/config"
	"url-shortener-pronetheus-consumer/internal/kafka"
	"url-shortener-pronetheus-consumer/internal/lib/logger/sl"
	"url-shortener-pronetheus-consumer/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(envLocal)

	c, err := kafka.NewConsumer(cfg, log, cfg.Topics...)
	if err != nil {
		log.Error("failed to create consumer", sl.Err(err))
	}

	metrics.Register()

	srv := http.NewServeMux()
	srv.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8090", srv); err != nil {
		log.Error("unable to start server", sl.Err(err))
	}

	signchan := make(chan os.Signal, 1)
	signal.Notify(signchan, syscall.SIGINT, syscall.SIGTERM)

	run := true 
	for run {
		select {
		case sig := <-signchan:
			log.Info(fmt.Sprintf("Caught signal: %v, terminating", sig))
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			log.Info(fmt.Sprintf("Consumed event from topic: %s: key=%s, value%s",
			*ev.TopicPartition.Topic, ev.Key, ev.Value))
			metrics.AuthCounter.WithLabelValues(string(ev.Key)).Inc()
		}
	}

	c.Close()
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
