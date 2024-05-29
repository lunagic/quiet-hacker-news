package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lunagic/environment-go/environment"
	"github.com/lunagic/quiet-hacker-news/internal/hackernews"
	"github.com/lunagic/quiet-hacker-news/internal/qhn"
)

type Config struct {
	Port int `env:"PORT"`
}

func main() {
	ctx := context.Background()

	config := Config{
		Port: 8080,
	}

	if err := environment.New().Decode(&config); err != nil {
		log.Fatalf("Error reading the environment: %s", err.Error())
	}

	s, err := qhn.New(
		hackernews.New(),
	)
	if err != nil {
		log.Fatalf("Error setting up service: %s", err.Error())
	}

	go keepRunning(time.Hour, func() {
		if err := s.Background(ctx); err != nil {
			log.Printf("Error running background job: %s", err.Error())
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Port),
		Handler: s,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}

func keepRunning(frequency time.Duration, action func()) {
	ticker := time.NewTicker(frequency)
	action()

	for range ticker.C {
		action()
	}
}
