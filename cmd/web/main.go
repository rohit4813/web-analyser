package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"web-analyser/config"
	"web-analyser/internal/router"
	"web-analyser/util/logger"
)

func main() {
	c := config.New()
	fmt.Println("debug level", c.Server.Debug)
	l := logger.New(c.Server.Debug)

	r := router.New(l)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.Server.Port),
		Handler: r,
	}

	closed := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals

		l.Info().Msgf("Shutting down server %v", s.Addr)
		close(closed)
	}()

	l.Info().Msgf("Starting server %v", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		l.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed
	l.Info().Msgf("Server shutdown successfully")
}
