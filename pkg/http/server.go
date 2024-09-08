package http

import (
	"context"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Run(ctx context.Context)
}

type api struct {
	router http.Handler
	port   string
}

func NewHTTPServer(router http.Handler, port string) Server {
	return &api{
		router: router,
		port:   port,
	}
}

func (a *api) Run(ctx context.Context) {
	a.serve(ctx)
}

// nolint:gosec
func (a *api) serve(ctx context.Context) {
	server := &http.Server{
		Addr:    ":" + a.port,
		Handler: a.router,
	}

	serverStopped := make(chan struct{})

	go func() {
		// Start the server in a separate goroutine
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Log if there's an error that isn't ErrServerClosed
			log.WithError(err).Error("Server ListenAndServe failed")
		}

		// Notify that the server has stopped
		close(serverStopped)
	}()

	// Log that the server is starting
	log.WithFields(log.Fields{"bind": a.port}).Info("Starting the API server")

	// Gracefully shut down when the context is done or the server stops
	select {
	case <-ctx.Done():
		log.Info("Shutting down the server due to context cancellation")

		// Shutdown the server with a background context
		if err := server.Shutdown(context.Background()); err != nil {
			log.WithError(err).Error("Server Shutdown failed")
		}
	case <-serverStopped:
		// The server stopped for some other reason (e.g., error)
		log.Info("Server stopped")
	}
}
