package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-analyser/api/backend/analyser"
	"web-analyser/api/router"
	"web-analyser/config"
	iHttp "web-analyser/internal/utils/http"
	"web-analyser/internal/utils/logger"
)

// setting up tpl global variable to load all the templates in memory to be rendered
var tpl *template.Template

// init loads all the templates in tpl
func init() {
	tpl = template.Must(template.ParseGlob("api/templates/*.gohtml"))
}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		panic("unable to load .env file")
	}

	// initialising the config
	conf := config.New()

	// setup required objects
	log := logger.NewLogger(conf.Server.Debug)
	httpClient := iHttp.NewHttpClient(&conf.Client)
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

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Hour)
	defer shutdownRelease()

	// serv.Shutdown(shutdownCtx) will gracefully shut down until the context is expired
	if err := serv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Shutdown error")
	}

	log.Info().Msg("Graceful shutdown complete")
}
