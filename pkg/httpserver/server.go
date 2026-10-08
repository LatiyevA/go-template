package httpserver

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

// Server owns the HTTP listener and its shutdown lifecycle.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	startOnce       sync.Once
}

// New configures a server. Call Start explicitly after initializing dependencies.
func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Handler:           handler,
			Addr:              ":8080",
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Start binds the configured address and reports immediate listener failures.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}
	return s.serve(listener)
}

func (s *Server) serve(listener net.Listener) error {
	started := false
	s.startOnce.Do(func() {
		started = true
		go func() {
			s.notify <- s.server.Serve(listener)
			close(s.notify)
		}()
	})
	if !started {
		if err := listener.Close(); err != nil {
			return err
		}
		return http.ErrServerClosed
	}
	return nil
}

// Notify reports when serving stops, including unexpected listener failures.
func (s *Server) Notify() <-chan error { return s.notify }

// Shutdown waits for in-flight requests until the configured timeout.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// Close terminates active connections after a failed graceful shutdown.
func (s *Server) Close() error { return s.server.Close() }
