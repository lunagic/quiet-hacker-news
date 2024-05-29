package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lunagic/quiet-hacker-news/internal/hackernews"
	"github.com/lunagic/quiet-hacker-news/internal/qhn"
)

func main() {
	ctx := context.Background()

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
		Addr:    fmt.Sprintf("0.0.0.0:%d", getPortNumber()),
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

func getPortNumber() int {
	defaultPort := 8080

	portString := os.Getenv("PORT")
	if portString == "" {
		return defaultPort
	}

	portInt, err := strconv.Atoi(portString)
	if err != nil {
		return defaultPort
	}

	return portInt
}
