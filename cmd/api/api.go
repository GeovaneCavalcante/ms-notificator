package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TIMEOUT = 30 * time.Second

type ServerOption func(server *http.Server)

func Start(port string, handler http.Handler, options ...ServerOption) error {
	srv := &http.Server{
		ReadTimeout:  TIMEOUT,
		WriteTimeout: TIMEOUT,
		Addr:         ":" + port,
		Handler:      handler,
	}

	for _, o := range options {
		o(srv)
	}

	var serverErr error

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		log.Println(fmt.Sprintf("Service listening on port %s", port))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("HTTP server error: %v", err)
			serverErr = err
			sigChannel <- syscall.SIGINT
		}
		log.Println("Stopped serving new connections")
	}()

	<-sigChannel

	log.Println("Stopping server")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), TIMEOUT)
	defer shutdownRelease()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		panic(err)
	}

	log.Println("Graceful shutdown complete")

	return serverErr
}

func WithReadTimeout(t time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.WriteTimeout = t
	}
}
