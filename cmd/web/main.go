package main

import (
	"fmt"
	"golang.org/x/net/html"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"web-analyser/config"
	"web-analyser/internal/router"
	"web-analyser/util/logger"
)

var tpl *template.Template

func print(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "title" {
		fmt.Printf("title:%+v\n", n.FirstChild.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		print(c)
	}
}

func main() {
	//url := "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\"><html><title></title></html>"
	//doc, _ := html.Parse(strings.NewReader(url))
	//print(doc)
	//publicKeyValue := "-//W3C//DTD HTML 3.2 Final//EN"
	//re := regexp.MustCompile("(^.*) ((xhtml|html) [0-9]*\\.?[0-9]*)(.*$)")
	//matches := re.FindStringSubmatch(strings.ToLower(publicKeyValue))
	//if len(matches) >= 3 {
	//	fmt.Println("version", matches[2])
	//}
	//parsedURL, err := url.ParseRequestURI("http://")
	//if err != nil {
	//	fmt.Println("error", err)
	//	return
	//}
	//
	//fmt.Println("host", parsedURL.Host)
	c := config.New()
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
