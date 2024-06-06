package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-analyser/api/resources/analyser"
	"web-analyser/config"
	"web-analyser/internal/router"
	h "web-analyser/util/http"
	"web-analyser/util/logger"
)

var tpl *template.Template

// init loads all the templates in tpl global variable
func init() {
	tpl = template.Must(template.ParseGlob("api/templates/*.gohtml"))
}

func main() {
	// initialising the config
	conf := config.New()
	log := logger.New(conf.Server.Debug)

	httpClient := h.NewClient(&conf.Client)
	a := analyser.NewAnalyser(httpClient)
	handler := analyser.NewHandler(log, tpl, a)
	mux := router.New(log, handler)

	serv := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Server.Port),
		ReadTimeout:  conf.Server.TimeoutRead,
		WriteTimeout: conf.Server.TimeoutWrite,
		Handler:      mux,
	}

	// Starting the server in another goroutine
	go func() {
		log.Info().Msgf("Starting server %v", serv.Addr)
		if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Server error")
		}
		log.Info().Msg("Stopped serving new connections")
	}()

	// listening for exit signals, this is a blocking operation
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	// serv.Shutdown(shutdownCtx) will gracefully shut down until the context is expired
	if err := serv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Shutdown error")
	}

	log.Info().Msg("Graceful shutdown complete")
}
